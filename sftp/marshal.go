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

func MarshalUint64(b []byte, u uint64) []byte {
	e := make([]byte, 8)

	binary.BigEndian.PutUint64(e, u)

	return append(b, e[0], e[1], e[2], e[3], e[4], e[5], e[6], e[7])
}

func UnmarshalUint64(b []byte) (uint64, []byte) {
	u := binary.BigEndian.Uint64(b)

	return u, b[8:]
}

func UnmarshalUint64Safe(b []byte) (uint64, []byte, error) {
	if len(b) < 8 {
		return 0, nil, errors.New("Not enough bytes to unmarshal uint64")
	}

	u, b := UnmarshalUint64(b)

	return u, b, nil
}

func MarshalString(b []byte, s string) []byte {
	return append(MarshalUint32(b, uint32(len(s))), s...)
}

func UnmarshalString(b []byte) (string, []byte) {
	s, b := UnmarshalUint32(b)

	return string(b[:s]), b[s:]
}

func UnmarshalStringSafe(b []byte) (string, []byte, error) {
	s, b, err := UnmarshalUint32Safe(b)

	if err != nil {
		return "", nil, err
	}

	if int64(s) > int64(len(b)) {
		return "", nil, errors.New("Not enough bytes to unmarshal string")
	}

	return string(b[:s]), b[s:], nil
}
