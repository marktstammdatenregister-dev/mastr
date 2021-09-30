package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/twpayne/go-geom"
	"github.com/twpayne/go-geom/encoding/ewkbhex"
	"github.com/twpayne/go-geom/encoding/wkt"
	//"github.com/twpayne/go-geom/xy"
	"io"
	"log"
	"os"
)

const squareMeterPerSquareDegree = float64(8500000000)

func convert(r *csv.Reader, w *csv.Writer, minArea float64) error {
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

		// Convert to MultiPolygon.
		mpoly := geom.NewMultiPolygon(geom.XY)
		switch t := g.(type) {
		case *geom.LineString:
			continue
			//lstr := geom.LineString(*t)
			//lrng := geom.NewLinearRingFlat(lstr.Layout(), lstr.FlatCoords())
			//if !xy.IsRingCounterClockwise(lrng.Layout(), lrng.FlatCoords()) {
			//	lrng.Reverse()
			//}
			//poly := geom.NewPolygon(geom.XY)
			//poly.Push(lrng)
			//mpoly.Push(poly)
		//case *geom.Polygon:
		//	poly := geom.Polygon(*t)
		//	mpoly.Push(&poly)
		case *geom.MultiPolygon:
			mpoly2 := geom.MultiPolygon(*t)
			mpoly = &mpoly2
		default:
			return fmt.Errorf("cannot handle geometry %v", t)
		}
		g = nil

		area := mpoly.Area()
		if area < 0 {
			return fmt.Errorf("negative area")
		}

		// Skip if area is too small.
		if minArea != 0 && area < minArea {
			continue
		}

		// Encode coordinates as WKT.
		h, err := wkt.Marshal(mpoly)
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
	r := csv.NewReader(bufio.NewReader(os.Stdin))
	r.Comma = '\t'
	r.LazyQuotes = true
	r.ReuseRecord = true

	// Construct CSV writer.
	bw := bufio.NewWriter(os.Stdout)
	defer bw.Flush()
	w := csv.NewWriter(bw)
	defer w.Flush()
	w.Comma = '\t'

	checkErr(convert(r, w, minAreaSquareDegrees))
}
