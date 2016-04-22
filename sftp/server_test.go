package sftp

import (
	"os"
	"sync"
	"testing"
)

func TestSendPacket(t *testing.T) {
	packet := &SSHFxInitPacket{
		Version: uint32(3),
	}

	server := &Server{
		In:       os.Stdin,
		Out:      os.Stdout,
		OutMutex: &sync.Mutex{},
	}

	n, err := server.SendPacket(packet)

	if err != nil {
		t.Error(err)
	}

	if uint32(n) != (packet.Length() + 4) {
		t.Errorf("Expected to send %v bytes, sent %v", packet.Length(), n)
	}
}
