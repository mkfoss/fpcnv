package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatalf("Usage: fpcnv [<OPEN_FILE>]")
	}

	fi, err := os.Stat(args[0])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println((float64(fi.Size()) / 1024.0) / 1024.0)

	//tst := &fpcnv.FPCnv{}
	//err = tst.Open(args[0])
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer func() { _ = tst.Close() }()
	//
	//fmt.Printf("header %#v\n\n", tst.Header)
	//
	//for idx := range tst.Fields.Count() {
	//	fld := tst.Fields.Field(idx)
	//	fmt.Printf("name: %s, type: %s, off: %d, size: %d, dec: %d \n", fld.Name(), fld.DbfType(), fld.Offset(),
	//		fld.Size(), fld.Decimals())
	//}
}
