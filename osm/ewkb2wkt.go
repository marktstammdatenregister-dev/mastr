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

		// Extract coordinates.
		var coords [][][]geom.Coord
		switch t := g.(type) {
		case *geom.LineString:
			g2 := geom.LineString(*t)
			coords = [][][]geom.Coord{[][]geom.Coord{g2.Coords()}}
		case *geom.Polygon:
			g2 := geom.Polygon(*t)
			coords = [][][]geom.Coord{g2.Coords()}
		case *geom.MultiPolygon:
			g2 := geom.MultiPolygon(*t)
			coords = g2.Coords()
		default:
			return fmt.Errorf("cannot handle geometry %v", t)
		}
		g = nil

		// Construct MultiPolygon.
		mpg := geom.NewMultiPolygon(geom.XY).MustSetCoords(coords)

		// Skip if area is too small.
		if minArea != 0 && mpg.Area() < minArea {
			continue
		}

		// Encode coordinates as WKT.
		h, err := wkt.Marshal(mpg)
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
