package fpcnv

import (
	"io"
	"os"
)

type FPFiler interface {
	io.Seeker
	io.Closer
	io.Reader
}

type FPOpener interface {
	OpenFile(filename string, flags int, perms os.FileMode) (FPFiler, error)
}

type OsOpener struct {
}

func (opn *OsOpener) OpenFile(filename string, flags int, perms os.FileMode) (FPFiler, error) {

	return os.OpenFile(filename, flags, perms)
}
