package sftp

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
)

// SFTP Packet Type Values
// https://tools.ietf.org/html/draft-ietf-secsh-filexfer-13#section-4.3
const (
	SSH_FXP_INIT           = 1
	SSH_FXP_VERSION        = 2
	SSH_FXP_OPEN           = 3
	SSH_FXP_CLOSE          = 4
	SSH_FXP_READ           = 5
	SSH_FXP_WRITE          = 6
	SSH_FXP_LSTAT          = 7
	SSH_FXP_FSTAT          = 8
	SSH_FXP_SETSTAT        = 9
	SSH_FXP_FSETSTAT       = 10
	SSH_FXP_OPENDIR        = 11
	SSH_FXP_READDIR        = 12
	SSH_FXP_REMOVE         = 13
	SSH_FXP_MKDIR          = 14
	SSH_FXP_RMDIR          = 15
	SSH_FXP_REALPATH       = 16
	SSH_FXP_STAT           = 17
	SSH_FXP_RENAME         = 18
	SSH_FXP_READLINK       = 19
	SSH_FXP_SYMLINK        = 20
	SSH_FXP_STATUS         = 101
	SSH_FXP_HANDLE         = 102
	SSH_FXP_DATA           = 103
	SSH_FXP_NAME           = 104
	SSH_FXP_ATTRS          = 105
	SSH_FXP_EXTENDED       = 200
	SSH_FXP_EXTENDED_REPLY = 201
)

type Packet interface {
	Type() int
	RequestId() uint32
	Length() uint32
	Payload() []byte
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	Response() (Packet, error)
}

// General Packet Format
// https://tools.ietf.org/html/draft-ietf-secsh-filexfer-13#section-4
func MarshalPacket(p Packet) ([]byte, error) {
	var encoded []byte

	encoded = MarshalUint32(encoded, p.Length())
	encoded = append(encoded, byte(p.Type()))

	if (p.Type() != SSH_FXP_INIT) && (p.Type() != SSH_FXP_VERSION) {
		encoded = MarshalUint32(encoded, p.RequestId())
	}

	encoded = append(encoded, p.Payload()...)

	return encoded, nil
}

func UnmarshalPacket(b []byte) (int, uint32, []byte, error) {
	var type_ int
	var requestId uint32
	var payload []byte
	var err error

	_, b, err = UnmarshalUint32Safe(b)

	if err != nil {
		return type_, requestId, payload, err
	}

	type_ = int(b[0])
	b = b[1:]

	if (type_ != SSH_FXP_INIT) && (type_ != SSH_FXP_VERSION) {
		requestId, b, err = UnmarshalUint32Safe(b)

		if err != nil {
			return type_, requestId, payload, err
		}
	}

	payload = b

	return type_, requestId, payload, err
}

func DecodePacket(b []byte) (Packet, error) {
	var packet Packet
	var err error

	type_ := int(b[4])

	switch type_ {
	case SSH_FXP_INIT:
		packet = &SSHFxInitPacket{}
	case SSH_FXP_VERSION:
		packet = &SSHFxVersionPacket{}
	case SSH_FXP_REALPATH:
		packet = &SSHFxRealPathPacket{}
	default:
		err = fmt.Errorf("Unrecognized packet type %v", type_)
	}

	if err != nil {
		log.WithFields(log.Fields{
			"type": type_,
		}).Error("Unrecognized packet type")

		return nil, err
	}

	err = packet.Unmarshal(b)

	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"type": packet.Type(),
	}).Debug("Decoded packet")

	return packet, err
}

type SSHFxInitPacket struct {
	Version uint32
}

func (p *SSHFxInitPacket) Type() int         { return SSH_FXP_INIT }
func (p *SSHFxInitPacket) RequestId() uint32 { return 0 }
func (p *SSHFxInitPacket) Length() uint32    { return 1 + 4 }

func (p *SSHFxInitPacket) Payload() []byte {
	var encoded []byte

	encoded = MarshalUint32(encoded, p.Version)

	return encoded
}

func (p *SSHFxInitPacket) Marshal() ([]byte, error) {
	encoded, err := MarshalPacket(p)

	return encoded, err
}

func (p *SSHFxInitPacket) Unmarshal(b []byte) error {
	_, _, payload, err := UnmarshalPacket(b)

	if err != nil {
		return err
	}

	version, payload, err := UnmarshalUint32Safe(payload)

	if err != nil {
		return err
	}

	p.Version = version

	return nil
}

func (p *SSHFxInitPacket) Response() (Packet, error) {
	packet := &SSHFxVersionPacket{
		Version: uint32(3),
	}

	return packet, nil
}

type SSHFxVersionPacket struct {
	Version    uint32
	Extensions []ExtensionPair
}

func (p *SSHFxVersionPacket) Type() int         { return SSH_FXP_VERSION }
func (p *SSHFxVersionPacket) RequestId() uint32 { return 0 }

func (p *SSHFxVersionPacket) Length() uint32 {
	length := 1 + 4

	for _, extension := range p.Extensions {
		// 4 bytes per string marshalled + length of each string
		length += 4 + len(extension.Name) + 4 + len(extension.Data)
	}

	return uint32(length)
}

func (p *SSHFxVersionPacket) Payload() []byte {
	var payload []byte

	payload = MarshalUint32(payload, p.Version)

	for _, extension := range p.Extensions {
		payload = MarshalString(payload, extension.Name)
		payload = MarshalString(payload, extension.Data)
	}

	return payload
}

func (p *SSHFxVersionPacket) Marshal() ([]byte, error) {
	encoded, err := MarshalPacket(p)

	return encoded, err
}

func (p *SSHFxVersionPacket) Unmarshal(b []byte) error {
	_, _, payload, err := UnmarshalPacket(b)

	if err != nil {
		return err
	}

	version, payload, err := UnmarshalUint32Safe(payload)

	if err != nil {
		return err
	}

	p.Version = version

	if len(payload) > 0 {
		var extension ExtensionPair

		extension, payload, err = UnmarshalExtensionPair(payload)

		if err != nil {
			return err
		}

		p.Extensions = append(p.Extensions, extension)
	}

	return nil
}

func (p *SSHFxVersionPacket) Response() (Packet, error) {
	return nil, nil
}

type SSHFxRealPathPacket struct {
	RequestID    uint32
	OriginalPath string
}

func (p *SSHFxRealPathPacket) Type() int         { return SSH_FXP_REALPATH }
func (p *SSHFxRealPathPacket) RequestId() uint32 { return p.RequestID }

func (p *SSHFxRealPathPacket) Length() uint32 {
	// Type + Request Id + String(4 + len(OriginalPath))
	length := 1 + 4 + 4 + len(p.OriginalPath)

	return uint32(length)
}

func (p *SSHFxRealPathPacket) Payload() []byte {
	var payload []byte

	payload = MarshalString(payload, p.OriginalPath)

	return payload
}

func (p *SSHFxRealPathPacket) Marshal() ([]byte, error) {
	encoded, err := MarshalPacket(p)

	return encoded, err
}

func (p *SSHFxRealPathPacket) Unmarshal(b []byte) error {
	_, requestId, payload, err := UnmarshalPacket(b)

	if err != nil {
		return err
	}

	originalPath, _, err := UnmarshalStringSafe(payload)

	if err != nil {
		return err
	}

	p.RequestID = requestId
	p.OriginalPath = originalPath

	return nil
}

func (p *SSHFxRealPathPacket) Response() (Packet, error) {
	// TODO: Implement SSH_FXP_NAME packet
	return nil, nil
}
