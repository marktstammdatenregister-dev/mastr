package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkbhex"
	"github.com/twpayne/go-geom/encoding/wkt"
	"io"
	"log"
	"os"
)

const squareMeterPerSquareDegree = float64(8500000000)

func convert(r *csv.Reader, w *csv.Writer, minArea float64) error {
	mps := 0
	other := 0
	for {
		// Read CSV record.
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Decode EWKB string.
		g, err := ewkbhex.Decode(record[0])
		if err != nil {
			return err
		}

		// Assert that we're dealing with a multipolygon.
		mp, ok := g.(*geom.MultiPolygon)
		if !ok {
			other += 1
			continue
		}
		mps += 1

		if minArea != 0 {
			area := mp.Area()
			if area < 0 {
				return fmt.Errorf("negative area")
			}

			// Skip if area is too small.
			if area < minArea {
				continue
			}
		}

		// Encode coordinates as WKT.
		h, err := wkt.Marshal(mp)
		if err != nil {
			return err
		}
		record[0] = h

		// Write CSV record.
		err = w.Write(record)
		if err != nil {
			return err
		}
	}
	fmt.Fprintf(os.Stderr, "ewkb2wkt: %d multipolygons, %d other (skipped)\n", mps, other)
	return nil
}

func checkErr(err error) {
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func main() {
	minAreaSquareMeters := flag.Float64("min-area", 0, "[square meters] buildings smaller than this are dropped")
	flag.Parse()

	if *minAreaSquareMeters < 0 {
		log.Fatal("min-area must not be negative")
	}
	minAreaSquareDegrees := *minAreaSquareMeters / squareMeterPerSquareDegree

	// Construct CSV reader.
	const bufSize = 4096 * 1024
	r := csv.NewReader(bufio.NewReaderSize(os.Stdin, bufSize))
	r.Comma = '\t'
	r.LazyQuotes = true
	r.ReuseRecord = true

	// Construct CSV writer.
	bw := bufio.NewWriterSize(os.Stdout, bufSize)
	defer bw.Flush()
	w := csv.NewWriter(bw)
	defer w.Flush()
	w.Comma = '\t'

	checkErr(convert(r, w, minAreaSquareDegrees))
}
