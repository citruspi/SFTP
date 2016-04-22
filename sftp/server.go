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
