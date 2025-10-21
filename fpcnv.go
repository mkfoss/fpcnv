package fpcnv

import (
	"io"
	"os"
)

type FPCnv struct {
	file   FPFiler
	Header *Header
	Fields *Fields
}

func (cnv *FPCnv) Open(filename string) error {
	return cnv.OpenWithOpener(filename, &OsOpener{})
}

func (cnv *FPCnv) OpenWithOpener(filename string, opener FPOpener) error {

	file, err := opener.OpenFile(filename, os.O_RDONLY, 0600)
	if err != nil {
		return NewErrorf("failed to open %s as a fp file", filename).SetWrapped(err)
	}

	cnv.Header, err = readHeader(file)
	if err != nil {
		return NewError("failed to read fp Header").SetWrapped(err)
	}

	cnv.Fields, err = readFields(file, int((cnv.Header.RecordsOffset-296)/32))
	if err != nil {
		return NewError("failed to read fp Fields").SetWrapped(err)
	}

	off, err := file.Seek(int64(cnv.Header.RecordsOffset), io.SeekStart)
	if err != nil {
		return NewError("failed to seek to start of records").SetWrapped(err)
	}

	if off != int64(cnv.Header.RecordsOffset) {
		return NewError("incorrect seek to end of records")
	}

	cnv.file = file

	return nil
}

func (cnv *FPCnv) Close() error {

	cnv.Header = nil
	cnv.Fields = nil
	err := cnv.file.Close()
	if err != nil {
		return NewErrorf("failed to close fp file").SetWrapped(err)
	}

	return nil
}
