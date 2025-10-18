package fpcnv

import (
	"encoding/binary"
	"slices"
)

type Header struct {
	Magic uint8 // file magic byte - what type of dbf
}

var mbsupported []byte = []byte{0x30}

func readHeader(f FPFiler) (*Header, error) {

	return nil, NewError("not implemented")
}

func readMagicByte(f FPFiler) (byte, error) {

	var b byte
	if err := binary.Read(f, binary.LittleEndian, &b); err != nil {
		return 0, NewError("failed to read magic byte").SetWrapped(err)
	}

	if !slices.Contains(mbsupported, b) {
		return 0, NewErrorf("unsupported dbf type")
	}

	return b, nil
}
