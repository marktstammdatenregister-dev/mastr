package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"database/sql"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/text/encoding/unicode"
	"io"
	"log"
	"os"
	"pvdb.de/mastr/internal"
	"pvdb.de/mastr/internal/spec"
	"strings"
	"text/template"
	"time"
)

var errMissingOption = errors.New("missing mandatory argument")

func main() {
	err := insert()
	if errors.Is(err, errMissingOption) {
		flag.PrintDefaults()
		os.Exit(64)
	}
	if err != nil {
		log.Fatalf("%v", err)
	}
}

func insert() error {
	const defaultOption = "<undefined>"
	exportFileName := flag.String("export", defaultOption, "file name of the export zip file")
	specFileName := flag.String("spec", defaultOption, "file name of the table spec")
	sqliteFile := flag.String("database", defaultOption, "file name of the SQLite database")
	flag.Parse()

	// Ensure mandatory flags are set.
	for _, arg := range []string{
		*exportFileName,
		*specFileName,
		*sqliteFile,
	} {
		if arg == defaultOption {
			return errMissingOption
		}
	}

	location, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		return fmt.Errorf("failed to load location data: %w", err)
	}
	internal.Location = location

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

	// Connect to the database.
	ctx := context.Background()
	db, err := sql.Open("sqlite3", *sqliteFile)
	if err != nil {
		return fmt.Errorf("failed to open: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("%v", err)
		}
	}()

	// Insert XML files one by one.
	dec := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	for _, ed := range export {
		createTable(db, ctx, &ed.Table)
		for _, xmlFile := range r.File {
			name := xmlFile.FileHeader.Name
			if !strings.HasPrefix(name, ed.Prefix) {
				continue
			}
			if err = func() error {
				fmt.Println(name)
				tx, err := db.BeginTx(ctx, nil) // non-nil TxOptions have no effect, see https://github.com/mattn/go-sqlite3/issues/685
				if err != nil {
					return fmt.Errorf("failed to start txn: %w", err)
				}

				//start := time.Now()
				f, err := xmlFile.Open()
				if err != nil {
					return fmt.Errorf("failed to open xml file %s: %w", name, err)
				}
				defer func() {
					if err := f.Close(); err != nil {
						log.Printf("%v", err)
					}
				}()
				_, err = insertFromXml(tx, ctx, dec.Reader(f), &ed.Table)
				if err != nil {
					return fmt.Errorf("failed to insert xml file %s: %w", name, err)
				}
				err = tx.Commit()
				if err != nil {
					return fmt.Errorf("failed to commit txn: %w", err)
				}
				//elapsed := time.Since(start).Seconds()
				//log.Printf("%s\t%.f entries/second", name, float64(i)/elapsed)
				return nil
			}(); err != nil {
				return fmt.Errorf("failed to process xml file: %w", err)
			}
		}
	}
	return nil
}

func createTable(db *sql.DB, ctx context.Context, td *spec.Table) error {
	// Generate "create table" statement.
	tmpl := template.Must(template.New("create").Parse(`
create table "{{.Element}}" (
	{{range .Fields -}}
		"{{.Name}}"
		{{- with .Sqlite}} {{.}}{{else}} text{{end}}
		{{- /* {{- if .Mandatory}} not null{{end}} ["mandatory" fields are frequently missing] */ -}}
		{{- with .References}} references "{{.Table}}"("{{.Column}}") deferrable initially deferred{{end}}
	{{end -}}
	primary key ("{{.Primary}}")
);
	`))
	var stmt bytes.Buffer
	if err := tmpl.Execute(&stmt, td); err != nil {
		return fmt.Errorf("failed to execute sql template: %w", err)
	}

	// Create the table.
	_, err := db.ExecContext(ctx, stmt.String())
	if err != nil {
		return fmt.Errorf("failed to execute create table statement: %w", err)
	}
	return nil
}

var prevPrimary string

func insertFromXml(tx *sql.Tx, ctx context.Context, f io.Reader, td *spec.Table) (int, error) {
	i := 0

	// Construct the buffered XML reader.
	const bufSize = 4096 * 1024
	br := xml.NewDecoder(bufio.NewReaderSize(f, bufSize))

	// Prepare the INSERT statement.
	fields := internal.NewFields(td.Fields)
	headers := fields.Header()
	columns := headers[0]
	placeholders := "?"
	for _, header := range headers[1:] {
		columns = fmt.Sprintf("%s, %s", columns, header)
		placeholders = fmt.Sprintf("%s, ?", placeholders)
	}

	// The export can contain duplicates, hence the UPSERT instead of a plain INSERT.
	stmt, err := tx.PrepareContext(ctx, fmt.Sprintf(`
insert into %s (%s) values (%s)
on conflict (%s) do nothing;
`, td.Element, columns, placeholders, td.Primary))
	if err != nil {
		return i, err
	}
	defer stmt.Close()

	// Copy data into the table.
	s := internal.NewXMLReader(td, br)
	for {
		values, err := s.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return i, err
		}
		record, err := fields.Record(values)
		if err != nil {
			return i, err
		}
		result, err := stmt.ExecContext(ctx, record...)
		if err != nil {
			return i, err
		}
		rows, err := result.RowsAffected()
		if err != nil {
			return i, err
		}
		if rows < 0 || rows > 1 {
			return i, fmt.Errorf("%d rows affected?!", rows)
		}
		if rows == 0 {
			fmt.Printf("- skipped duplicate: %s.%s=%s\n", td.Element, td.Primary, values[td.Primary])
		}
		i++
	}
	return i, nil
}
