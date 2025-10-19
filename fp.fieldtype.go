package fpcnv

import "bytes"

type FieldType byte

const (
	FTUnknown FieldType = iota
	FTCharacter
	FTCurrency
	FTNumeric
	FTFloat
	FTDate
	FTDateTime
	FTDouble
	FTInteger
	FTLogical
	FTMemo
	FTGeneral
	FTPicture
)

var supportedfieldtypes = []byte("CYNFDTBILMGP")

func (f FieldType) String() string {
	if f < 1 || int(f) > len(supportedfieldtypes) {
		return "unknown"
	}
	return string(supportedfieldtypes[int(f)-1])
}

func FieldTypeFromByte(ftchar byte) FieldType {

	idx := bytes.IndexByte([]byte(supportedfieldtypes), bytes.ToUpper([]byte{ftchar})[0])
	if idx == -1 {
		return FTUnknown
	}
	return FieldType(idx + 1)
}
