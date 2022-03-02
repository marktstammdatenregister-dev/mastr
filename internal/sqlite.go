package internal

import (
	"bytes"
	"fmt"
	"io"
	"text/template"
	"time"

	"crawshaw.io/sqlite"
	"crawshaw.io/sqlite/sqlitex"
	"marktstammdatenregister.dev/internal/spec"
)

type SqliteWriter struct {
	pool *sqlitex.Pool

	// Per-table state.
	conn   *sqlite.Conn
	td     *spec.Table
	query  string
	fields *Fields
}

var _ Recorder = (*SqliteWriter)(nil)
var _ io.Closer = (*SqliteWriter)(nil)

func NewSqliteWriter(db string) (*SqliteWriter, error) {
	location, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		return nil, fmt.Errorf("failed to load location data: %w", err)
	}
	Location = location

	pool, err := sqlitex.Open(db, 0, 1)
	if err != nil {
		return nil, err
	}

	conn := pool.Get(nil)
	defer pool.Put(conn)

	// Single writer, no reader, don't care about integrity if the import fails.
	// This roughly halves the runtime of the SQLite writer.
	sqlitex.Exec(conn, "PRAGMA journal_mode=MEMORY", nil)
	sqlitex.Exec(conn, "PRAGMA synchronous=OFF", nil)
	sqlitex.Exec(conn, "PRAGMA locking_mode=EXCLUSIVE", nil)

	return &SqliteWriter{
		pool:  pool,
		query: "not a valid SQL query!",
	}, nil
}

// EnterTable implements Recorder.
func (w *SqliteWriter) EnterTable(td spec.Table) error {
	w.td = &td

	// Generate "create table" statement.
	tmpl := template.Must(template.New("create").Parse(`
create table "{{.Element}}" (
	{{range .Fields -}}
		"{{.Name}}"
		{{- with .Sqlite}} {{.}}{{else}} text{{end}}
		{{- with .References}} references "{{.Table}}"("{{.Column}}"){{end}},
	{{end -}}
	primary key ("{{.Primary}}")
)
	`))
	var stmt bytes.Buffer
	if err := tmpl.Execute(&stmt, td); err != nil {
		return fmt.Errorf("failed to execute sql template: %w", err)
	}

	w.conn = w.pool.Get(nil)
	if err := sqlitex.Exec(w.conn, stmt.String(), nil); err != nil {
		return fmt.Errorf("failed to execute create table statement: %w", err)
	}

	w.fields = NewFields(td.Fields)
	headers := w.fields.Header()
	columns := headers[0]
	placeholders := "?"
	for _, header := range headers[1:] {
		columns = fmt.Sprintf("%s, %s", columns, header)
		placeholders = fmt.Sprintf("%s, ?", placeholders)
	}

	// The export can contain duplicates, hence the UPSERT instead of a plain INSERT.
	w.query = fmt.Sprintf(
		`insert into %s (%s) values (%s) on conflict (%s) do nothing;`,
		td.Element, columns, placeholders, td.Primary)

	return nil
}

// LeaveTable implements Recorder.
func (w *SqliteWriter) LeaveTable() error {
	// TODO: Add sanity checks.

	// Create an index for each field with the "Index" flag.
	table := w.td.Element
	for _, field := range w.td.Fields {
		if !field.Index {
			continue
		}
		col := field.Name
		stmt := fmt.Sprintf(`create index "idx_%s_%s" on "%s"("%s")`, table, col, table, col)
		if err := sqlitex.Exec(w.conn, stmt, nil); err != nil {
			return err
		}
	}

	w.td = nil
	w.pool.Put(w.conn)
	w.conn = nil
	return nil
}

// EnterFile implements Recorder.
func (w *SqliteWriter) EnterFile(f string) error {
	return nil
}

// LeaveFile implements Recorder.
func (w *SqliteWriter) LeaveFile() error {
	return nil
}

// Record implements Recorder.
func (w *SqliteWriter) Record(item map[string]string) error {
	rec, err := w.fields.Record(item)
	if err != nil {
		return fmt.Errorf("failed to write record: %w", err)
	}
	return sqlitex.Exec(w.conn, w.query, nil, rec...)
}

// Close implements io.Closer.
func (w *SqliteWriter) Close() error {
	if err := func() error {
		conn := w.pool.Get(nil)
		defer w.pool.Put(conn)
		if err := sqlitex.Exec(conn, "analyze", nil); err != nil {
			return err
		}
		return sqlitex.Exec(conn, "vacuum", nil)
	}(); err != nil {
		return fmt.Errorf("failed to vacuum: %s", err)
	}

	return w.pool.Close()
}
