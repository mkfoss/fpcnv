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

	if err := cnv.OpenWithFiler(file); err != nil {
		return NewErrorf("failed to open %s as a fp file", filename).SetWrapped(err)
	}

	return nil
}

func (cnv *FPCnv) OpenWithFiler(filer FPFiler) error {

	if _, err := filer.Seek(0, io.SeekStart); err != nil {
		return NewErrorf("failed to seek to begin of file at start").SetWrapped(err)
	}

	var err error
	cnv.Header, err = ReadHeader(filer)
	if err != nil {
		return NewError("failed to read fp Header").SetWrapped(err)
	}

	cnv.Fields, err = ReadFields(filer, int((cnv.Header.RecordsOffset-296)/32))
	if err != nil {
		return NewError("failed to read fp Fields").SetWrapped(err)
	}

	off, err := filer.Seek(int64(cnv.Header.RecordsOffset), io.SeekStart)
	if err != nil {
		return NewError("failed to seek to start of records").SetWrapped(err)
	}

	if off != int64(cnv.Header.RecordsOffset) {
		return NewError("incorrect seek to end of records")
	}

	cnv.file = filer

	return nil
}

func (cnv *FPCnv) Close() error {

	cnv.Header = nil
	cnv.Fields = nil
	err := cnv.file.Close()
	if err != nil {
		return NewErrorf("failed to close fp file").SetWrapped(err)
	}
	cnv.file = nil
	return nil
}

func (cnv *FPCnv) Active() bool {

	return cnv.file != nil
}

func (cnv *FPCnv) InitializeProcessor(processor RecordProcessor) error {
	if !cnv.Active() {
		return NewErrorf("cannot process records if no file is open")
	}
	err := processor.Initialize(cnv.file, cnv.Header, cnv.Fields)
	if err != nil {
		return NewErrorf("failed to initialize processor").SetWrapped(err)
	}
	return nil
}
