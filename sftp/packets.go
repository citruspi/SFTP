package sftp

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
}

// General Packet Format
// https://tools.ietf.org/html/draft-ietf-secsh-filexfer-13#section-4
func MarshalPacket(type_ int, requestID uint32, payload []byte) []byte {
	var encoded []byte

	encoded = []byte{0, 0, 0, 0}

	encoded = append(encoded, byte(type_))

	if (type_ != SSH_FXP_INIT) && (type_ != SSH_FXP_VERSION) {
		encoded = MarshalUint32(encoded, requestID)
	}

	encoded = append(encoded, payload...)

	length := []byte{}

	length = MarshalUint32(length, uint32(len(encoded)-4))

	encoded[0], encoded[1], encoded[2], encoded[3] = length[0], length[1], length[2], length[3]

	return encoded
}

type SSHFxInitPacket struct {
	Version uint32
}

func (p SSHFxInitPacket) MarshalBinary() ([]byte, error) {
	var encoded []byte

	encoded = MarshalUint32(encoded, p.Version)

	encoded = MarshalPacket(SSH_FXP_INIT, 0, encoded)

	return encoded, nil
}

func (p *SSHFxInitPacket) UnmarshalBinary(b []byte) error {
	var version uint32
	var err error

	version, b, err = UnmarshalUint32Safe(b[1:])

	if err != nil {
		return err
	}

	p.Version = version

	return nil
}

type SSHFxVersionPacket struct {
	Version    uint32
	Extensions []ExtensionPair
}

func (p SSHFxVersionPacket) MarshalBinary() ([]byte, error) {
	// 1 byte for the packet type and 4 for the version value
	length := 1 + 4

	for _, extension := range p.Extensions {
		// 4 bytes per string marshalled + length of each string
		length += 4 + len(extension.Name) + 4 + len(extension.Data)
	}

	binary := make([]byte, 0, length)

	binary = append(binary, SSH_FXP_INIT)
	binary = MarshalUint32(binary, p.Version)

	for _, extension := range p.Extensions {
		binary = MarshalString(binary, extension.Name)
		binary = MarshalString(binary, extension.Data)
	}

	return binary, nil
}
