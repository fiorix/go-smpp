package smpp34

import (
	"encoding/binary"
)

func unpackUi32(b []byte) (n uint32) {
	n = binary.BigEndian.Uint32(b)
	return
}

func packUi32(n uint32) (b []byte) {
	b = make([]byte, 4)
	binary.BigEndian.PutUint32(b, n)
	return
}

func unpackUi16(b []byte) (n uint16) {
	n = binary.BigEndian.Uint16(b)
	return
}

func packUi16(n uint16) (b []byte) {
	b = make([]byte, 2)
	binary.BigEndian.PutUint16(b, n)
	return
}
