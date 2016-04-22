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

func TestSSHFxVersionPacketMarshal(t *testing.T) {
	packet := &SSHFxVersionPacket{
		Version:    uint32(3),
		Extensions: []ExtensionPair{},
	}

	marshalled, err := packet.Marshal()

	if err != nil {
		t.Error(err)
	}

	expected := []byte{0, 0, 0, 5, 2, 0, 0, 0, 3}

	if !bytes.Equal(marshalled, expected) {
		t.Errorf("Expected %v, received %v", expected, marshalled)
	}
}
