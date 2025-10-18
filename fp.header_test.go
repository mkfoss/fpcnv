package fpcnv

import (
	"errors"
	"fmt"
	"testing"

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

			mb, err := readMagicByte(fl)
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

		mb, err := readMagicByte(fl)
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
