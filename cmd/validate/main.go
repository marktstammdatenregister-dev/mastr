package main

import (
	"archive/zip"
	"bufio"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"golang.org/x/text/encoding/unicode"
	"io"
	"log"
	"marktstammdatenregister.dev/internal"
	"marktstammdatenregister.dev/internal/spec"
	"os"
	"strings"
)

var errMissingOption = errors.New("missing mandatory argument")

func main() {
	err := validate()
	if errors.Is(err, errMissingOption) {
		flag.PrintDefaults()
		os.Exit(64)
	}
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func validate() error {
	const defaultOption = "<undefined>"
	exportFileName := flag.String("export", defaultOption, "file name of the export zip file")
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

	v := internal.NewValidator(
		strings.SplitN(*exportFileName, "_", 2)[0],
		fmt.Sprintf("https://download.marktstammdatenregister.de/%s", *exportFileName),
		os.Stderr,
		os.Stdout,
	)
	dec := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	for _, ed := range export {
		if err := func() error {
			if err := v.EnterTable(ed.Table); err != nil {
				return fmt.Errorf("failed to enter table: %s", err)
			}
			defer func() {
				if err := v.LeaveTable(); err != nil {
					log.Printf("%v", err)
				}
			}()

			for _, xmlFile := range r.File {
				name := xmlFile.FileHeader.Name
				if !strings.HasPrefix(name, ed.Prefix) {
					continue
				}
				if err = func() error {
					if err := v.EnterFile(name); err != nil {
						return fmt.Errorf("failed to enter file: %s", err)
					}
					defer func() {
						if err := v.LeaveFile(); err != nil {
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
					_, err = processFile(v, dec.Reader(f), &ed.Table)
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
	if err := v.Close(); err != nil {
		return fmt.Errorf("failed to close validator: %w", err)
	}
	return nil
}

func processFile(rec internal.Recorder, f io.Reader, td *spec.Table) (int, error) {
	// Construct the buffered XML reader.
	const bufSize = 4096 * 1024
	d := xml.NewDecoder(bufio.NewReaderSize(f, bufSize))
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

		if err := rec.Record(item); err != nil {
			return i, err
		}
	}
	return i, nil
}
