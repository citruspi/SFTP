package sftp

import (
	"io"
	"sync"

	log "github.com/Sirupsen/logrus"
)

type Server struct {
	In            io.Reader
	Out           io.WriteCloser
	OutMutex      *sync.Mutex
	PacketChannel chan Packet
	WorkerCount   int
}

func (s *Server) SendPacket(p Packet) (int, error) {
	log.WithFields(log.Fields{
		"type": p.Type(),
		"id":   p.RequestId(),
	}).Debug("Marshalling packet")

	encoded, err := p.Marshal()

	if err != nil {
		return 0, err
	}

	s.OutMutex.Lock()
	defer s.OutMutex.Unlock()

	log.WithFields(log.Fields{
		"type": p.Type(),
		"id":   p.RequestId(),
	}).Debug("Sending packet")

	n, err := s.Out.Write(encoded)

	return n, err
}

func (s *Server) ReceivePackets() error {
	defer close(s.PacketChannel)

	for {
		var encoded []byte
		var length_encoded []byte
		var body []byte

		length_encoded = make([]byte, 4)

		s.In.Read(length_encoded)

		length, _, err := UnmarshalUint32Safe(length_encoded)

		if err != nil {
			return err
		}

		body = make([]byte, length)

		_, err = io.ReadFull(s.In, body)

		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		encoded = append(length_encoded, body...)

		packet, err := DecodePacket(encoded)

		if err != nil {
			return err
		}

		log.WithFields(log.Fields{
			"type": packet.Type(),
			"id":   packet.RequestId(),
		}).Debug("Received packet")

		s.PacketChannel <- packet
	}
}

func (s *Server) Worker(results chan error) {
	for packet := range s.PacketChannel {
		response, err := packet.Response()

		if err != nil {
			results <- err
			return
		}

		log.WithFields(log.Fields{
			"packet_type":   packet.Type(),
			"packet_id":     packet.RequestId(),
			"response_type": response.Type(),
			"response_id":   response.RequestId(),
		}).Debug("Responding to packet")

		if response != nil {
			_, err = s.SendPacket(response)

			if err != nil {
				results <- err
				return
			}
		}
	}

	results <- nil
}

func (s *Server) Serve() error {
	log.Debug("Starting SFTP server")

	defer s.Out.Close()

	log.Debug("Receiving packets")

	go s.ReceivePackets()

	results := make(chan error)

	log.Debug("Starting server workers")
	for i := 0; i < s.WorkerCount; i++ {
		go s.Worker(results)
	}

	log.Debug("Checking for worker errors")
	for i := 0; i < s.WorkerCount; i++ {
		err := <-results

		if err != nil {
			return err
		}
	}

	log.Debug("Stopping SFTP server")
	return nil
}

func NewServer(in io.Reader, out io.WriteCloser, workers int) (*Server, error) {
	server := &Server{
		In:            in,
		Out:           out,
		OutMutex:      &sync.Mutex{},
		PacketChannel: make(chan Packet),
		WorkerCount:   workers,
	}

	log.WithFields(log.Fields{
		"workers": workers,
	}).Info("Created new server")

	return server, nil
}
