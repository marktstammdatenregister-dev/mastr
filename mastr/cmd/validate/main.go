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
	"os"
	"pvdb.de/mastr/internal"
	"pvdb.de/mastr/internal/spec"
	"sort"
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

	v := &validator{
		key:      make(map[string]map[string]int),
		files:    make([]string, 0),
		errCount: 0,
	}
	dec := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	for _, ed := range export {
		for _, xmlFile := range r.File {
			name := xmlFile.FileHeader.Name
			if !strings.HasPrefix(name, ed.Prefix) {
				continue
			}
			if err = func() error {
				f, err := xmlFile.Open()
				if err != nil {
					return fmt.Errorf("failed to open xml file %s: %w", name, err)
				}
				defer func() {
					if err := f.Close(); err != nil {
						log.Printf("%v", err)
					}
				}()
				_, err = v.validateFile(dec.Reader(f), name, &ed.Table)
				if err != nil {
					return fmt.Errorf("failed to validate xml file %s: %w", name, err)
				}
				return nil
			}(); err != nil {
				return fmt.Errorf("failed to process xml file: %w", err)
			}
		}
	}
	if v.errCount == 0 {
		println("SUCCESS")
	} else {
		fmt.Printf("FAILURE: %d error(s) found\n", v.errCount)
	}
	return nil
}

type duplicate struct {
	table         string
	column        string
	key           string
	originalFile  string
	duplicateFile string
}

type broken struct {
	table         string
	column        string
	key           string
	primary       string
	targetTable   string
	targetColumn  string
	targetKey     string
	referenceFile string
}

type missing struct {
	firstKey string
	count    int
}

type validator struct {
	key      map[string]map[string]int
	files    []string
	errCount int
}

func (v *validator) validateFile(f io.Reader, fileName string, td *spec.Table) (int, error) {
	fmt.Println(fileName)
	// Insert the file name into `v.files`.
	for _, name := range v.files {
		if name == fileName {
			return 0, fmt.Errorf("file %s already validated", fileName)
		}
	}
	v.files = append(v.files, fileName)
	fileIndex := len(v.files) - 1

	// Construct the buffered XML reader.
	const bufSize = 4096 * 1024
	d := xml.NewDecoder(bufio.NewReaderSize(f, bufSize))
	r := internal.NewXMLReader(td, d)

	// Validate the file.
	mis := make(map[string]missing)
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

		// Check for duplicate key definitions.
		key := item[td.Primary]
		if _, ok := v.key[td.Element]; !ok {
			v.key[td.Element] = make(map[string]int)
		}
		keys := v.key[td.Element]
		if originalFileIndex, ok := keys[key]; ok {
			reportDuplicate(duplicate{
				table:         td.Element,
				column:        td.Primary,
				key:           key,
				originalFile:  v.files[originalFileIndex],
				duplicateFile: fileName,
			})
			v.errCount++
		}
		keys[key] = fileIndex

		// Check for mandatory fields. Reported once per file.
		for _, field := range td.Fields {
			if !field.Mandatory {
				continue
			}
			if _, ok := item[field.Name]; !ok {
				if _, ok := mis[field.Name]; !ok {
					mis[field.Name] = missing{firstKey: item[td.Primary], count: 0}
				}
				m := mis[field.Name]
				m.count++
				mis[field.Name] = m
				v.errCount++
			}
		}

		// Check for broken references.
		for _, field := range td.Fields {
			ref := field.References
			if ref == nil {
				continue
			}
			if _, ok := item[field.Name]; !ok {
				continue
			}
			x := item[field.Name]
			brk := broken{
				table:         td.Element,
				column:        field.Name,
				key:           key,
				primary:       td.Primary,
				targetTable:   ref.Table,
				targetColumn:  ref.Column,
				targetKey:     x,
				referenceFile: fileName,
			}
			if _, ok := v.key[ref.Table]; !ok {
				reportBroken(brk)
				v.errCount++
				continue
			}
			if _, ok := v.key[ref.Table][x]; !ok {
				reportBroken(brk)
				v.errCount++
				continue
			}
		}
	}
	reportMissing(mis, td.Element, td.Primary)
	return i, nil
}

func reportDuplicate(dup duplicate) {
	fmt.Printf("- duplicate: %s.%s=%s already appeared in %s\n", dup.table, dup.column, dup.key, dup.originalFile)
}

func reportBroken(brk broken) {
	fmt.Printf("- broken: %s(%s=%s).%s references %s.%s=%s, which is missing\n", brk.table, brk.primary, brk.key, brk.column, brk.targetTable, brk.targetColumn, brk.targetKey)
}

func reportMissing(mis map[string]missing, table string, primary string) {
	cols := make([]string, 0)
	for col, _ := range mis {
		cols = append(cols, col)
	}
	sort.Strings(cols)
	for _, col := range cols {
		m := mis[col]
		fmt.Printf("- missing: %s.%s is mandatory but missing (%d times, e.g. %s=%s)\n", table, col, m.count, primary, m.firstKey)
	}
}
