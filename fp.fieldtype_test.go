package fpcnv

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FieldTypeString(t *testing.T) {

	assert.Equal(t, FTUnknown, FieldTypeFromByte(byte(0xFF)))
	assert.Equal(t, FTUnknown, FieldTypeFromByte(byte(0x00)))
	assert.Equal(t, FTUnknown, FieldTypeFromByte(byte(0x12)))

	assert.Equal(t, FTUnknown, FieldTypeFromByte([]byte("X")[0]))
	assert.Equal(t, FTCharacter, FieldTypeFromByte([]byte("C")[0]))
	assert.Equal(t, FTCurrency, FieldTypeFromByte([]byte("Y")[0]))
	assert.Equal(t, FTNumeric, FieldTypeFromByte([]byte("N")[0]))
	assert.Equal(t, FTFloat, FieldTypeFromByte([]byte("F")[0]))
	assert.Equal(t, FTDate, FieldTypeFromByte([]byte("D")[0]))
	assert.Equal(t, FTDateTime, FieldTypeFromByte([]byte("T")[0]))
	assert.Equal(t, FTDouble, FieldTypeFromByte([]byte("B")[0]))
	assert.Equal(t, FTInteger, FieldTypeFromByte([]byte("I")[0]))
	assert.Equal(t, FTLogical, FieldTypeFromByte([]byte("L")[0]))
	assert.Equal(t, FTMemo, FieldTypeFromByte([]byte("M")[0]))
	assert.Equal(t, FTGeneral, FieldTypeFromByte([]byte("G")[0]))
	assert.Equal(t, FTPicture, FieldTypeFromByte([]byte("P")[0]))

	assert.Equal(t, FTUnknown, FieldTypeFromByte([]byte("x")[0]))
	assert.Equal(t, FTCharacter, FieldTypeFromByte([]byte("c")[0]))
	assert.Equal(t, FTCurrency, FieldTypeFromByte([]byte("y")[0]))
	assert.Equal(t, FTNumeric, FieldTypeFromByte([]byte("n")[0]))
	assert.Equal(t, FTFloat, FieldTypeFromByte([]byte("f")[0]))
	assert.Equal(t, FTDate, FieldTypeFromByte([]byte("d")[0]))
	assert.Equal(t, FTDateTime, FieldTypeFromByte([]byte("t")[0]))
	assert.Equal(t, FTDouble, FieldTypeFromByte([]byte("b")[0]))
	assert.Equal(t, FTInteger, FieldTypeFromByte([]byte("i")[0]))
	assert.Equal(t, FTLogical, FieldTypeFromByte([]byte("l")[0]))
	assert.Equal(t, FTMemo, FieldTypeFromByte([]byte("m")[0]))
	assert.Equal(t, FTGeneral, FieldTypeFromByte([]byte("g")[0]))
	assert.Equal(t, FTPicture, FieldTypeFromByte([]byte("p")[0]))

	for idx, _ := range supportedfieldtypes {
		typ := FieldType(idx + 1)
		assert.Equal(t, typ.String(), string(supportedfieldtypes[int(typ-1)]))
		assert.Equal(t, FieldType(typ), FieldTypeFromByte(supportedfieldtypes[idx]))
	}
}
