package fpcnv

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testField(t *testing.T, fld Fielder) {
}

func Test_BaseField(t *testing.T) {

	type testtype struct {
		id       int
		input    []byte
		expected *BaseField
		error    string
	}

	cases := []testtype{
		{1, []byte{0x41, 0x47, 0x45, 0x4E, 0x43, 0x59, 0x49, 0x44, 0x00, 0x00, 0x00, 0x43, 0x01, 0x00, 0x00, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			&BaseField{
				name:     "agencyid",
				dbft:     FTCharacter,
				offset:   1,
				size:     5,
				decimals: 0,
				system:   false,
				nullable: false,
				binary:   false,
			},
			""},
		{2, []byte{0x4D, 0x4F, 0x44, 0x49, 0x46, 0x49, 0x45, 0x44, 0x4F, 0x4E, 0x00, 0x44, 0xC1, 0x02, 0x00, 0x00, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			&BaseField{
				name:     "modifiedon",
				dbft:     FTDate,
				offset:   705,
				size:     8,
				decimals: 0,
				system:   false,
				nullable: false,
				binary:   false,
			},
			"",
		},
		{3, []byte{0x47, 0x52, 0x4F, 0x53, 0x53, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x59, 0x41, 0x00, 0x00, 0x00, 0x08, 0x04, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			&BaseField{
				name:     "gross",
				dbft:     FTCurrency,
				offset:   65,
				size:     8,
				decimals: 4,
				system:   false,
				nullable: false,
				binary:   true,
			},
			"",
		},
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("Field read #%.2d", tc.id), func(t *testing.T) {
			fl := &MockFiler{
				Data: tc.input,
			}
			bf, err := ReadField(fl)
			if tc.error != "" {
				if err == nil {
					t.Errorf("ReadField(%d) expected error: %s", tc.id, tc.error)
				}
				assert.ErrorContains(t, err, tc.error)
				assert.Nil(t, bf)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, bf)
			}
		})
	}
}
