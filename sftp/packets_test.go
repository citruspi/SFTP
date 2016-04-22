package sftp

import (
	"bytes"
	"testing"
)

func TestSSHFxInitPacketMarshalBinary(t *testing.T) {
	var version uint32
	var extensions []ExtensionPair

	version = uint32(3)

	packet := &SSHFxInitPacket{
		Version:    version,
		Extensions: extensions,
	}

	marshalled, err := packet.MarshalBinary()

	if err != nil {
		t.Error(err)
	}

	expected := []byte{1, 0, 0, 0, 3}

	if !bytes.Equal(marshalled, expected) {
		t.Errorf("Expected %v, received %v", expected, marshalled)
	}
}
