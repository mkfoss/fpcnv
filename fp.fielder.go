package fpcnv

import (
	"bytes"
	"encoding/binary"
)

type Fielder interface {
	Name() string
	DbfType() FieldType
	Offset() uint32
	Size() uint8
	Decimals() uint8
	System() bool
	Nullable() bool
	Binary() bool
	Value() any
	String() string
}

type BaseField struct {
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

func ReadField(filer FPFiler) (*BaseField, error) {
	bf := fieldReadStruct{}
	err := binary.Read(filer, binary.LittleEndian, &bf)
	if err != nil {
		return nil, NewError("reading field information").SetWrapped(err)
	}
	return &BaseField{
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

func (bf *BaseField) Name() string {
	return bf.name
}

func (bf *BaseField) DbfType() FieldType {
	return bf.dbft
}

func (bf *BaseField) Offset() uint32 {
	return bf.offset
}

func (bf *BaseField) Size() uint8 {
	return bf.size
}

func (bf *BaseField) Decimals() uint8 {
	return bf.decimals
}

func (bf *BaseField) System() bool {
	return bf.system
}

func (bf *BaseField) Nullable() bool {
	return bf.nullable
}

func (bf *BaseField) Binary() bool {
	return bf.Binary()
}

type Fields struct {
	fields   []Fielder
	fieldmap map[string]int
}

func ReadFields(filer FPFiler) (*Fields, error) {

	return nil, NewError("not implemented")
}
