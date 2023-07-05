package main

import (
	"archive/zip"
	"bufio"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/unicode"
	"io"
	"io/ioutil"
	"log"
	"marktstammdatenregister.dev/internal"
	"marktstammdatenregister.dev/internal/spec"
	"os"
	"strings"
)

var errMissingOption = errors.New("missing mandatory argument")

func main() {
	err := mainWithError()
	if errors.Is(err, errMissingOption) {
		flag.PrintDefaults()
		os.Exit(64)
	}
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func mainWithError() error {
	const defaultOption = "<undefined>"
	sqliteFile := flag.String("database", defaultOption, "(optional) file name of the SQLite database")
	exportFileName := flag.String("export", defaultOption, "file name of the export zip file")
	textReport := flag.String("report-text", "stderr", "(optional) text report output [stdout | stderr | off]")
	jsonReport := flag.String("report-json", "off", "(optional) json report output [stdout | stderr | off]")
	specFileName := flag.String("spec", defaultOption, "file name of the export spec")
	flag.Parse()

	// Ensure mandatory flags are set.
	for _, arg := range []string{
		*exportFileName,
		*specFileName,
	} {
		if arg == defaultOption {
			return errMissingOption
		}
	}

	// Instantiate report writers.
	writers := map[string]io.Writer{
		"stdout": os.Stdout,
		"stderr": os.Stderr,
		"off":    ioutil.Discard,
	}
	var textWriter, jsonWriter io.Writer
	textWriter = os.Stderr
	jsonWriter = ioutil.Discard
	if *textReport != defaultOption {
		w, ok := writers[*textReport]
		if !ok {
			return fmt.Errorf("invalid output: %s", *textReport)
		}
		textWriter = w
	}
	if *jsonReport != defaultOption {
		w, ok := writers[*jsonReport]
		if !ok {
			return fmt.Errorf("invalid output: %s", *jsonReport)
		}
		jsonWriter = w
	}

	export, err := spec.DecodeExport(*specFileName)
	if err != nil {
		return fmt.Errorf("failed to decode export spec: %w", err)
	}

	r, err := zip.OpenReader(*exportFileName)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer func() {
		if err := r.Close(); err != nil {
			log.Printf("%v", err)
		}
	}()

	exportName, _, found := strings.Cut(*exportFileName, "_")
	if !found {
		return fmt.Errorf("export file name '%s' does not contain an underscore", *exportFileName)
	}
	recs := []internal.Recorder{
		internal.NewUnusedTracker(r.File, textWriter),
		internal.NewValidator(
			exportName,
			fmt.Sprintf("https://download.marktstammdatenregister.de/%s", *exportFileName),
			textWriter,
			jsonWriter,
		)}
	if *sqliteFile != defaultOption {
		w, err := internal.NewSqliteWriter(*sqliteFile)
		if err != nil {
			return err
		}
		recs = append(recs, w)
	}
	dec := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	for _, ed := range export {
		if err := func() error {
			if err := enterTable(recs, ed.Table); err != nil {
				return fmt.Errorf("failed to enter table: %s", err)
			}
			defer func() {
				if err := leaveTable(recs); err != nil {
					log.Printf("%v", err)
				}
			}()

			for _, xmlFile := range r.File {
				name := xmlFile.FileHeader.Name
				if !strings.HasPrefix(name, ed.Prefix) {
					continue
				}
				if err = func() error {
					if err := enterFile(recs, name); err != nil {
						return fmt.Errorf("failed to enter file: %s", err)
					}
					defer func() {
						if err := leaveFile(recs); err != nil {
							log.Printf("%v", err)
						}
					}()

					f, err := xmlFile.Open()
					if err != nil {
						return fmt.Errorf("failed to open xml file %s: %w", name, err)
					}
					defer func() {
						if err := f.Close(); err != nil {
							log.Printf("%v", err)
						}
					}()
					_, err = processFile(recs, dec.Reader(f), &ed.Table)
					if err != nil {
						return fmt.Errorf("failed to validate xml file %s: %w", name, err)
					}
					return nil
				}(); err != nil {
					return fmt.Errorf("failed to process xml file: %w", err)
				}
			}
			return nil
		}(); err != nil {
			return fmt.Errorf("failed to process xml file: %w", err)
		}
	}
	for _, rec := range recs {
		if err := rec.Close(); err != nil {
			return fmt.Errorf("failed to close recorder: %w", err)
		}
	}
	return nil
}

func enterTable(recs []internal.Recorder, td spec.Table) error {
	for _, rec := range recs {
		if err := rec.EnterTable(td); err != nil {
			return err
		}
	}
	return nil
}

func leaveTable(recs []internal.Recorder) error {
	for _, rec := range recs {
		if err := rec.LeaveTable(); err != nil {
			return err
		}
	}
	return nil
}

func enterFile(recs []internal.Recorder, f string) error {
	for _, rec := range recs {
		if err := rec.EnterFile(f); err != nil {
			return err
		}
	}
	return nil
}

func leaveFile(recs []internal.Recorder) error {
	for _, rec := range recs {
		if err := rec.LeaveFile(); err != nil {
			return err
		}
	}
	return nil
}

func processFile(recs []internal.Recorder, f io.Reader, td *spec.Table) (int, error) {
	// Construct the buffered XML reader.
	const bufSize = 4096 * 1024
	d := xml.NewDecoder(bufio.NewReaderSize(f, bufSize))
	d.CharsetReader = charset.NewReaderLabel
	r := internal.NewXMLReader(td, d)

	// Validate the file.
	i := 0
	for {
		item, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return i, err
		}
		i++

		for _, rec := range recs {
			if err := rec.Record(item); err != nil {
				return i, err
			}
		}
	}
	return i, nil
}
