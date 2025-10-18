package fpcnv

import (
	"encoding/binary"
	"slices"
	"time"
)

type Header struct {
	Magic         uint8 // file magic byte - what type of dbf
	LastUpdate    time.Time
	RecordCount   uint32
	RecordSize    uint16
	RecordsOffset uint16
	HasIndex      bool
	HasFpt        bool
	IsDatabase    bool
	Codepage      Codepage
}

var mbsupported = []byte{0x30}

func readHeader(f FPFiler) (*Header, error) {

	return nil, NewError("not implemented")
}

func readHdrMagicByte(f FPFiler) (byte, error) {
	var b byte
	if err := binary.Read(f, binary.LittleEndian, &b); err != nil {
		return 0, NewError("failed to read magic byte").SetWrapped(err)
	}

	if !slices.Contains(mbsupported, b) {
		return 0, NewErrorf("unsupported dbf type")
	}

	return b, nil
}

func readHdrLastUpdate(f FPFiler) (time.Time, error) {

	b := make([]byte, 3)
	if err := binary.Read(f, binary.LittleEndian, &b); err != nil {
		return time.Time{}, NewError("failed to read last update time").SetWrapped(err)
	}

	errformatstr := "invalid %s: %d in last updated date"

	if b[0] > 99 {
		return time.Time{}, NewErrorf(errformatstr, "year", b[0])
	}
	yr := int(b[0]) //year 200 compensation
	if b[0] >= 70 {
		yr += 1900
	} else {
		yr += 2000
	}

	if b[1] == 0 || b[1] > 12 {
		return time.Time{}, NewErrorf(errformatstr, "month", b[1])
	}

	daysinmnth := time.Date(yr, time.Month(b[1]+1), 0, 0, 0, 0, 0, time.Local).Day()
	if b[2] == 0 || int(b[2]) > daysinmnth {
		return time.Time{}, NewErrorf(errformatstr, "day", b[2])
	}

	return time.Date(yr, time.Month(b[1]), int(b[2]), 0, 0, 0, 0, time.Local), nil //time for dbf files always local
}

func readHdrNumRecords(f FPFiler) (uint32, error) {

	var b uint32
	err := binary.Read(f, binary.LittleEndian, &b)
	if err != nil {
		return 0, NewError("failed to read num records").SetWrapped(err)
	}
	if b == 0 {
		return 0, NewErrorf("invalid record count, dbf files require at least one record").SetWrapped(nil)
	}
	return b, nil
}

func readUint16(f FPFiler) (uint16, error) {

	var b uint16
	err := binary.Read(f, binary.LittleEndian, &b)
	if err != nil {
		return 0, NewError("failed to read uint16").SetWrapped(err)
	}

	return b, nil
}

func readHdrRecordOffset(f FPFiler) (uint16, error) {

	b, err := readUint16(f)
	if err != nil {
		return 0, NewError("failed to read offset").SetWrapped(err)
	}

	if b < 33 {
		return 0, NewErrorf("invalid record offset").SetWrapped(nil)
	}

	return b, nil
}

func readHdrRecordSize(f FPFiler) (uint16, error) {

	b, err := readUint16(f)
	if err != nil {
		return 0, NewError("failed to read record size").SetWrapped(err)
	}

	if b < 2 {
		return 0, NewErrorf("invalid record size").SetWrapped(nil)
	}

	return b, nil
}

func readHdrTableFlags(f FPFiler) (bool, bool, error) {

	var b byte
	if err := binary.Read(f, binary.LittleEndian, &b); err != nil {
		return false, false, NewError("failed to read flags").SetWrapped(err)
	}

	if b&0x04 == 0x04 {
		return false, false, NewErrorf("fp databases not supported").SetWrapped(nil)
	}

	if b > 7 {
		return false, false, NewErrorf("invalid table flags").SetWrapped(nil)
	}

	return b&0x01 == 0x01, b&0x02 == 0x02, nil
}

func readHdrCodepage(f FPFiler) (Codepage, error) {

	var b uint8
	if err := binary.Read(f, binary.LittleEndian, &b); err != nil {
		return 0, NewError("failed to read codepage").SetWrapped(err)
	}

	if !slices.Contains(supportedCodepages, Codepage(b)) {
		return 0, NewErrorf("unsupported codepage 0x%.2x", b)
	}

	return Codepage(b), nil
}
