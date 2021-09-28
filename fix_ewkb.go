package main

import (
	"encoding/binary"
	"encoding/csv"
	"github.com/twpayne/go-geom/encoding/ewkbhex"
	"io"
	"log"
	"os"
)

func convert(r *csv.Reader, w *csv.Writer) error {
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		g, err := ewkbhex.Decode(record[0])
		if err != nil {
			return err
		}
		h, err := ewkbhex.Encode(g, binary.LittleEndian)
		if err != nil {
			return err
		}
		record[0] = h
		err = w.Write(record)
		if err != nil {
			return err
		}
	}
	w.Flush()
	return nil
}

func main() {
	r := csv.NewReader(os.Stdin)
	r.Comma = '\t'
	r.LazyQuotes = true
	r.ReuseRecord = true

	w := csv.NewWriter(os.Stdout)
	w.Comma = '\t'

	err := convert(r, w)
	if err != nil {
		log.Fatalf("%v", err)
	}
}
