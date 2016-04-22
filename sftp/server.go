package sftp

import (
	"io"
	"sync"
)

type Server struct {
	In       io.Reader
	Out      io.WriteCloser
	OutMutex *sync.Mutex
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
