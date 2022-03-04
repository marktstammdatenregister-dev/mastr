package internal

import (
	"archive/zip"
	"fmt"
	"io"
	"sort"

	"marktstammdatenregister.dev/internal/spec"
)

type UnusedTracker struct {
	textWriter io.Writer

	currentTable     string
	currentTableUsed bool

	unusedTds     []string
	unparsedFiles map[string]struct{}
}

var _ Recorder = (*UnusedTracker)(nil)

func NewUnusedTracker(files []*zip.File, textWriter io.Writer) *UnusedTracker {
	unparsedFiles := make(map[string]struct{})
	for _, f := range files {
		unparsedFiles[f.FileHeader.Name] = struct{}{}
	}
	return &UnusedTracker{
		textWriter: textWriter,

		currentTable:     "",
		currentTableUsed: false,

		unusedTds:     make([]string, 0),
		unparsedFiles: unparsedFiles,
	}
}

// EnterTable implements Recorder.
func (w *UnusedTracker) EnterTable(td spec.Table) error {
	w.currentTable = td.Element
	w.currentTableUsed = false
	return nil
}

// LeaveTable implements Recorder.
func (w *UnusedTracker) LeaveTable() error {
	if !w.currentTableUsed {
		w.unusedTds = append(w.unusedTds, w.currentTable)
	}
	return nil
}

// EnterFile implements Recorder.
func (w *UnusedTracker) EnterFile(f string) error {
	w.currentTableUsed = true
	delete(w.unparsedFiles, f)
	return nil
}

// LeaveFile implements Recorder.
func (w *UnusedTracker) LeaveFile() error {
	return nil
}

// Record implements Recorder.
func (w *UnusedTracker) Record(item map[string]string) error {
	return nil
}

// Close implements io.Closer.
func (w *UnusedTracker) Close() error {
	// We're only reporting the last error that occurred. Fine for now.
	//var err error
	for _, t := range w.unusedTds {
		fmt.Fprintf(w.textWriter, "- unused: table spec %s\n", t)
		//err = fmt.Errorf("table spec %s not used, the export may be incomplete", t)
	}
	for _, f := range w.unparsedFilesSorted() {
		fmt.Fprintf(w.textWriter, "- unparsed: file %s\n", f)
		//err = fmt.Errorf("file %s not parsed, the spec may be incomplete", f)
	}
	//return err
	return nil
}

func (w *UnusedTracker) unparsedFilesSorted() []string {
	ret := make([]string, 0)
	for f, _ := range w.unparsedFiles {
		ret = append(ret, f)
	}
	sort.Strings(ret)
	return ret
}
