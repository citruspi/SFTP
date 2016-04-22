package sftp

import (
	"encoding/binary"
	"errors"
)

func MarshalUint32(b []byte, u uint32) []byte {
	e := make([]byte, 4)

	binary.BigEndian.PutUint32(e, u)

	return append(b, e[0], e[1], e[2], e[3])
}

func UnmarshalUint32(b []byte) (uint32, []byte) {
	u := binary.BigEndian.Uint32(b)

	return u, b[4:]
}

func UnmarshalUint32Safe(b []byte) (uint32, []byte, error) {
	if len(b) < 4 {
		return 0, nil, errors.New("Not enough bytes to unmarshal uint32")
	}

	u, b := UnmarshalUint32(b)

	return u, b, nil
}
