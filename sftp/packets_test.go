package sftp

import (
	"bytes"
	"reflect"
	"testing"
)

func TestSSHFxInitPacketMarshal(t *testing.T) {
	packet := &SSHFxInitPacket{
		Version: uint32(3),
	}

	marshalled, err := packet.Marshal()

	if err != nil {
		t.Error(err)
	}

	expected := []byte{0, 0, 0, 5, 1, 0, 0, 0, 3}

	if !bytes.Equal(marshalled, expected) {
		t.Errorf("Expected %v, received %v", expected, marshalled)
	}
}

func TestSSHFxInitPacketUnmarshal(t *testing.T) {
	packet := &SSHFxInitPacket{}

	err := packet.Unmarshal([]byte{0, 0, 0, 5, 1, 0, 0, 0, 3})

	if err != nil {
		t.Error(err)
	}

	expected := &SSHFxInitPacket{
		Version: 3,
	}

	if !reflect.DeepEqual(expected, packet) {
		t.Errorf("Expected %+v, received %+v", expected, packet)
	}
}

func TestSSHFxVersionPacketMarshalBinary(t *testing.T) {
	var version uint32
	var extensions []ExtensionPair

	version = uint32(3)

	packet := &SSHFxVersionPacket{
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
