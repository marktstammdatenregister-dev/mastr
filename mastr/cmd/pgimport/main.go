package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"context"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"github.com/jackc/pgx/v4"
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
	filePrefix := flag.String("prefix", defaultOption, "prefix of xml files to extract")
	databaseUrl := flag.String("database", defaultOption, "postgres database URL")
	forceCreate := flag.Bool("force-create", false, "use CREATE instead of CREATE IF NOT EXISTS")
	flag.Parse()

	// Ensure mandatory flags are set.
	for _, arg := range []string{
		*exportFileName,
		*specFileName,
		*filePrefix,
		*databaseUrl,
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

	td, err := spec.Decode(*specFileName)
	if err != nil {
		return fmt.Errorf("failed to decode spec: %w", err)
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
	conn, err := pgx.Connect(ctx, *databaseUrl)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer func() {
		if err := conn.Close(ctx); err != nil {
			log.Printf("%v", err)
		}
	}()

	// Insert XML files one by one.
	dec := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	for _, xmlFile := range r.File {
		if !strings.HasPrefix(xmlFile.FileHeader.Name, *filePrefix) {
			continue
		}
		if err = func() error {
			start := time.Now()
			f, err := xmlFile.Open()
			if err != nil {
				return fmt.Errorf("failed to open xml file: %w", err)
			}
			defer func() {
				if err := f.Close(); err != nil {
					log.Printf("%v", err)
				}
			}()
			i, err := insertFromXml(dec.Reader(f), conn, ctx, td, *forceCreate)
			if err != nil {
				return fmt.Errorf("failed to insert from xml file: %w", err)
			}
			elapsed := time.Since(start).Seconds()
			log.Printf("%s\t%.f entries/second", xmlFile.FileHeader.Name, float64(i)/elapsed)
			return nil
		}(); err != nil {
			return fmt.Errorf("failed to process xml file: %w", err)
		}
	}
	return nil
}

func createTable(tx pgx.Tx, ctx context.Context, td *spec.Table, force bool) error {
	// Generate "create table" statement.
	tmpl := template.Must(template.New("create").Parse(`
create unlogged table {{if .Force}}{{else}}if not exists{{end}}
{{- with .Descriptor}}"{{.Element}}" (
	{{range .Fields -}}
		"{{.Name}}"
		{{- with .Psql}} {{.}}{{else}} text{{end}}
		{{- if .Mandatory}} not null{{end}}
		{{- with .References}} references "{{.Table}}"("{{.Column}}") deferrable initially deferred{{end}},
	{{end -}}
	primary key ("{{.Primary}}")
) with (autovacuum_enabled=false);{{end}}
	`))
	var stmt bytes.Buffer
	if err := tmpl.Execute(&stmt, struct {
		Force      bool
		Descriptor *spec.Table
	}{force, td}); err != nil {
		return fmt.Errorf("failed to execute sql template: %w", err)
	}

	// Create the table.
	_, err := tx.Exec(ctx, stmt.String())
	if err != nil {
		return fmt.Errorf("failed to execute create table statement: %w", err)
	}
	return nil
}

func insertFromXml(f io.Reader, conn *pgx.Conn, ctx context.Context, td *spec.Table, force bool) (int64, error) {
	// Construct the buffered XML reader.
	const bufSize = 4096 * 1024
	br := xml.NewDecoder(bufio.NewReaderSize(f, bufSize))

	// Begin the transaction.
	tx, err := conn.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer func() {
		err := tx.Rollback(ctx)
		if err != nil && err != pgx.ErrTxClosed {
			log.Printf("%v", err)
		}
	}()

	// Create the table.
	if err := createTable(tx, ctx, td, force); err != nil {
		return 0, err
	}

	// Copy data into the table.
	fields := internal.NewFields(td.Fields)
	s := internal.NewXMLSource(td, br, fields)
	i, err := tx.CopyFrom(
		ctx,
		pgx.Identifier{td.Element},
		fields.Header(),
		&s,
	)
	if err != nil {
		return i, err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return i, err
	}
	return i, nil
}
