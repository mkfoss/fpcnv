package fpcnv

import (
	"bytes"
	"encoding/binary"
)

type Field struct {
	name     string
	dbft     FieldType
	offset   uint32
	size     uint8
	decimals uint8
	system   bool
	nullable bool
	binary   bool
}

type fieldReadStruct struct {
	Name      [11]byte
	Dbft      uint8
	Offset    uint32
	Size      uint8
	Decimals  uint8
	Flags     byte
	Reserved1 [13]byte
}

func readField(filer FPFiler) (*Field, error) {
	bf := fieldReadStruct{}
	err := binary.Read(filer, binary.LittleEndian, &bf)
	if err != nil {
		return nil, NewError("reading field information").SetWrapped(err)
	}

	return &Field{
		name:     string(bytes.ToLower(bytes.Trim(bf.Name[:], string([]byte{0x00})))),
		dbft:     FieldTypeFromByte(bf.Dbft),
		offset:   bf.Offset,
		size:     bf.Size,
		decimals: bf.Decimals,
		system:   bf.Flags&0x01 == 001,
		nullable: bf.Flags&0x02 == 0x02,
		binary:   bf.Flags&0x04 == 0x04,
	}, nil
}

func (bf *Field) Name() string {
	return bf.name
}

func (bf *Field) DbfType() FieldType {
	return bf.dbft
}

func (bf *Field) Offset() uint32 {
	return bf.offset
}

func (bf *Field) Size() uint8 {
	return bf.size
}

func (bf *Field) Decimals() uint8 {
	return bf.decimals
}

func (bf *Field) System() bool {
	return bf.system
}

func (bf *Field) Nullable() bool {
	return bf.nullable
}

func (bf *Field) Binary() bool {
	return bf.Binary()
}
