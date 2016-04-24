package sftp

import (
	"io"
	"sync"
)

type Server struct {
	In            io.Reader
	Out           io.WriteCloser
	OutMutex      *sync.Mutex
	PacketChannel chan Packet
	WorkerCount   int
}

func (s *Server) SendPacket(p Packet) (int, error) {
	encoded, err := p.Marshal()

	if err != nil {
		return 0, err
	}

	s.OutMutex.Lock()
	defer s.OutMutex.Unlock()

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
	defer s.Out.Close()

	go s.ReceivePackets()

	results := make(chan error)

	for i := 0; i < s.WorkerCount; i++ {
		go s.Worker(results)
	}

	for i := 0; i < s.WorkerCount; i++ {
		err := <-results

		if err != nil {
			return err
		}
	}

	return nil
}
