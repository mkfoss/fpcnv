package fpcnv

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_HeaderReadMagicbyte(t *testing.T) {

	type testType struct {
		id       int
		input    uint8
		expected uint8
		error    string
	}

	cases := []testType{
		{1, 0x30, 0x30, ""},
		{2, 0x31, 0x00, "unsupported dbf type"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("Magic Byte #%.2d", tc.id), func(t *testing.T) {

			fl := NewMockFiler()
			fl.Data = make([]byte, 32)
			fl.Data[0] = tc.input

			mb, err := readHdrMagicByte(fl)
			if tc.error != "" {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				assert.ErrorContains(t, err, tc.error)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, mb)
			}
		})
	}

	t.Run("Magic Byte - No data", func(t *testing.T) {
		fl := NewMockFiler()

		mb, err := readHdrMagicByte(fl)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "failed to read magic byte")

		var goterr Errorer
		ok := errors.As(err, &goterr)
		assert.True(t, ok)
		assert.Error(t, goterr)
		assert.NotNil(t, goterr.Unwrap())
		assert.ErrorContains(t, goterr.Unwrap(), "EOF")

		assert.Empty(t, mb)
	})
}

func Test_HeaderReadLastModified(t *testing.T) {

	type testtype struct {
		id       int
		input    []byte
		expected string
		error    string
	}

	cases := []testtype{
		{1, []byte{25, 10, 10}, "2025-10-10", ""},
		{2, []byte{32, 2, 5}, "2032-02-05", ""},
		{3, []byte{89, 10, 11}, "1989-10-11", ""},
		{4, []byte{100, 10, 11}, "0001-01-01", "invalid year: 100 in last updated date"},
		{5, []byte{23, 13, 11}, "0001-01-01", "invalid month: 13 in last updated date"},
		{6, []byte{25, 10, 32}, "0001-01-01", "invalid day: 32 in last updated date"},
		{7, []byte{20, 2, 29}, "2020-02-29", ""},
		{8, []byte{20, 2, 30}, "0001-01-01", "invalid day: 30 in last updated date"},
		{9, []byte{19, 2, 29}, "0001-01-01", "invalid day: 29 in last updated date"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("Read LastUpdated #%.2d", tc.id), func(t *testing.T) {

			expectedtm, err := time.ParseInLocation("2006-01-02", tc.expected, time.Local)
			if err != nil {
				t.Fatal(err, "should not fail")
			}
			fl := NewMockFiler()
			fl.Data = make([]byte, 32)
			fl.Data[1] = tc.input[0]
			fl.Data[2] = tc.input[1]
			fl.Data[3] = tc.input[2]
			fl.Pos = 1

			gottm, err := readHdrLastUpdate(fl)
			if tc.error != "" {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				assert.ErrorContains(t, err, tc.error)
				assert.Equal(t, gottm, time.Time{})
			} else {
				assert.NoError(t, err)
				assert.Equal(t, expectedtm, gottm)
			}
		})
	}

}

func Test_HeaderReadNumRecords(t *testing.T) {
	type testType struct {
		id       int
		input    []byte
		expected uint32
		error    string
	}

	cases := []testType{
		{1, []byte{0x0c, 0x10, 0x00, 0x00}, uint32(4108), ""},
		{2, []byte{0x00, 0x00, 0x00, 0x00}, uint32(0), "invalid record count, dbf files require at least one record"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("Read NumRecords #%.2d", tc.id), func(t *testing.T) {
			fl := NewMockFiler()
			fl.Data = make([]byte, 32)
			fl.Data[4] = tc.input[0]
			fl.Data[5] = tc.input[1]
			fl.Data[6] = tc.input[2]
			fl.Data[7] = tc.input[3]
			fl.Pos = 4

			gottm, err := readHdrNumRecords(fl)
			if tc.error != "" {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				assert.ErrorContains(t, err, tc.error)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, gottm)
			}
		})
	}
}

func Test_ReadUint16(t *testing.T) {
	type testType struct {
		id       int
		input    []byte
		expected uint16
		error    string
	}

	cases := []testType{
		{1, []byte{0x0c, 0x10}, uint16(4108), ""},
		{2, []byte{0x00, 0x00}, uint16(0), ""},
		{3, []byte{0x88, 0x04}, uint16(1160), ""},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("Read Uint16 #%.2d", tc.id), func(t *testing.T) {
			fl := NewMockFiler()
			fl.Data = make([]byte, 32)
			fl.Data[8] = tc.input[0]
			fl.Data[9] = tc.input[1]
			fl.Pos = 8

			gotint, err := readUint16(fl)
			if tc.error != "" {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				assert.ErrorContains(t, err, tc.error)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, gotint)
			}
		})
	}
}

func Test_ReadRecOffset(t *testing.T) {

	type testType struct {
		id       int
		input    []byte
		expected uint16
		error    string
	}

	cases := []testType{
		{1, []byte{0x0c, 0x10}, uint16(4108), ""},
		{2, []byte{0x00, 0x00}, uint16(0), "invalid record offset"},
		{3, []byte{0x01, 0x00}, uint16(0), "invalid record offset"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("Read RecOffset #%.2d", tc.id), func(t *testing.T) {
			fl := NewMockFiler()
			fl.Data = make([]byte, 32)
			fl.Data[10] = tc.input[0]
			fl.Data[11] = tc.input[1]
			fl.Pos = 10

			ro, err := readHdrRecordOffset(fl)
			if tc.error != "" {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				assert.ErrorContains(t, err, tc.error)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, ro)
			}
		})
	}
}

func Test_ReadRecordSize(t *testing.T) {
	type testType struct {
		id       int
		input    []byte
		expected uint16
		error    string
	}
	cases := []testType{
		{1, []byte{0x0c, 0x10}, uint16(4108), ""},
		{2, []byte{0x00, 0x00}, uint16(0), "invalid record size"},
		{3, []byte{0x01, 0x00}, uint16(0), "invalid record size"},
		{3, []byte{0x02, 0x00}, uint16(2), ""},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("Read RecSize #%.2d", tc.id), func(t *testing.T) {
			fl := NewMockFiler()
			fl.Data = make([]byte, 32)
			fl.Data[10] = tc.input[0]
			fl.Data[11] = tc.input[1]
			fl.Pos = 10

			b, err := readHdrRecordSize(fl)
			if tc.error != "" {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				assert.ErrorContains(t, err, tc.error)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, uint16(b))
			}
		})
	}
}

func Test_ReadTableFlags(t *testing.T) {
	type testType struct {
		id     int
		input  byte
		hasidx bool
		hasfpt bool
		error  string
	}

	cases := []testType{
		{1, 0, false, false, ""},
		{2, 1, true, false, ""},
		{3, 2, false, true, ""},
		{4, 3, true, true, ""},
		{5, 4, false, false, "fp databases not supported"},
		{6, 5, true, false, "fp databases not supported"},
		{7, 6, false, true, "fp databases not supported"},
		{8, 7, true, true, "fp databases not supported"},
		{9, 8, false, false, "invalid table flags"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("Read TableFlags #%.2d", tc.id), func(t *testing.T) {
			fl := NewMockFiler()
			fl.Data = make([]byte, 32)
			fl.Data[28] = tc.input
			fl.Pos = 28

			hasidx, hsfpt, err := readHdrTableFlags(fl)
			if tc.error != "" {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				assert.ErrorContains(t, err, tc.error)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.hasidx, hasidx)
				assert.Equal(t, tc.hasfpt, hsfpt)
			}
		})
	}
}

func Test_ReadCodepage(t *testing.T) {
	type testType struct {
		id       int
		input    byte
		expected Codepage
		error    string
	}

	cases := []testType{
		{1, 0x03, Codepage(0x03), ""},
		{2, 0x00, Codepage(0x00), "unsupported codepage 0x00"},
		{3, 0xff, Codepage(0x00), "unsupported codepage 0xff"},
		{3, 0x4b, Codepage(0x00), "unsupported codepage 0x4b"},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("Read Codepage #%.2d", tc.id), func(t *testing.T) {
			fl := NewMockFiler()
			fl.Data = make([]byte, 32)
			fl.Data[29] = tc.input
			fl.Pos = 29

			cp, err := readHdrCodepage(fl)
			if tc.error != "" {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				assert.ErrorContains(t, err, tc.error)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, cp)
			}
		})
	}
}
