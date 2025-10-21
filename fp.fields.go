package fpcnv

type Fields struct {
	fields   []*Field
	fieldmap map[string]int
}

// ReadFields reads the Fields from the file data - filer must be positioned correctly.
func readFields(filer FPFiler, expected int) (*Fields, error) {

	flds := &Fields{
		fields:   make([]*Field, expected),
		fieldmap: make(map[string]int),
	}

	var curroffset uint32 = 1 //offset starts at 1 - deleted marker

	for x := range expected {
		fld, err := readField(filer)
		if err != nil {
			return nil, NewErrorf("reading field %d of %d", x, expected).SetWrapped(err)
		}

		if fld.offset == 0 {
			fld.offset = curroffset
		}
		curroffset += uint32(fld.size)

		flds.fields[x] = fld
		flds.fieldmap[fld.name] = x
	}

	var b []byte = make([]byte, 1)
	n, err := filer.Read(b)
	if err != nil {
		return nil, NewError("reading end of Fields marker").SetWrapped(err)
	}
	if n != 1 {
		return nil, NewError("reading end of Fields marker, short read").SetWrapped(err)
	}
	if b[0] != 0x0D {
		return nil, NewError("incorrects end of Fields marker read")
	}

	return flds, nil
}

func (flds *Fields) Count() int {

	return len(flds.fields)
}

func (flds *Fields) Field(i int) *Field {
	if i >= 0 && i < len(flds.fields) {
		return flds.fields[i]
	}
	return nil
}

func (flds *Fields) FieldByName(name string) *Field {
	idx, ok := flds.fieldmap[name]
	if !ok {
		return nil
	}
	return flds.Field(idx)
}
