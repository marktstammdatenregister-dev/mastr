package internal

import (
	"fmt"
	"sort"

	"marktstammdatenregister.dev/internal/spec"
)

// JSON reporting.
type ExportReport struct {
	ExportName string
	Url        string
	Files      []FileReport
}

type FileReport struct {
	FileName   string
	NumMissing int
	NumBroken  int
	Missing    []MissingReport
	Broken     []BrokenReport
}

type MissingReport struct {
	FieldName      string
	Example        string
	NumOccurrences int
}

type BrokenReport struct {
	SourceKey       string
	ForeignKeyField string
	TargetTable     string
	TargetKeyField  string
	TargetKey       string
}

type Validator struct {
	key      map[string]map[string]int
	files    []string
	errCount int
	report   ExportReport

	// Table state.
	td *spec.Table

	// File state.
	fileState *fileState
}

type fileState struct {
	index   int
	report  FileReport
	missing map[string]missing
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

func (v *Validator) EnterTable(td spec.Table) {
	v.td = &td
}
func (v *Validator) LeaveTable() {
	v.td = nil
}
func (v *Validator) EnterFile(f string) {
	// Sanity check: have we seen this file before?
	for _, name := range v.files {
		if name == f {
			panic(fmt.Sprintf("file %s already validated", f))
		}
	}

	// Insert the file name into `v.files`.
	v.files = append(v.files, f)

	// Initialize file specific state.
	v.fileState = &fileState{
		report: FileReport{
			FileName: f,
			Missing:  make([]MissingReport, 0),
			Broken:   make([]BrokenReport, 0),
		},
		index:   len(v.files) - 1,
		missing: make(map[string]missing),
	}
}
func (v *Validator) LeaveFile() {
	s := v.fileState

	// Report missing files.
	cols := make([]string, 0)
	for col, _ := range s.missing {
		cols = append(cols, col)
	}
	sort.Strings(cols)

	for _, col := range cols {
		m := s.missing[col]
		s.report.Missing = append(s.report.Missing, MissingReport{
			FieldName:      col,
			Example:        m.firstKey,
			NumOccurrences: m.count,
		})
		s.report.NumMissing += m.count
	}

	// Append this file report to the validation report.
	v.report.Files = append(v.report.Files, s.report)
}
func (v *Validator) Record(r map[string]string) {
	// TODO ...


	// Check for duplicate key definitions.
	key := item[td.Primary]
	if _, ok := v.key[td.Element]; !ok {
		v.key[td.Element] = make(map[string]int)
	}
	keys := v.key[td.Element]
	if originalFileIndex, ok := keys[key]; ok {
		if outputText {
			reportDuplicate(duplicate{
				table:         td.Element,
				column:        td.Primary,
				key:           key,
				originalFile:  v.files[originalFileIndex],
				duplicateFile: fileName,
			})
		}
		v.errCount++

		// TODO: Include duplicates in JSON.
	}
	keys[key] = fileIndex

}
