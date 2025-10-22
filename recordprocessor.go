package fpcnv

type RecordProcessor interface {
	Initialize(filer FPFiler, header *Header, fields *Fields) error
}
