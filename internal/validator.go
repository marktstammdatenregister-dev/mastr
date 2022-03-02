package internal

import (
	"encoding/json"
	"fmt"
	"io"

	"marktstammdatenregister.dev/internal/spec"
)

type Validator struct {
	textWriter io.Writer
	jsonWriter io.Writer

	key      map[string]map[string]int
	files    []string
	errCount int
	report   ExportReport

	// Table state.
	td *spec.Table

	// File state.
	fileState *fileState
}

var _ Recorder = (*Validator)(nil)
var _ io.Closer = (*Validator)(nil)

// JSON reporting.
type ExportReport struct {
	ExportName string
	Url        string
	Files      []FileReport
}

type FileReport struct {
	FileName   string
	NumBroken  int
	Broken     []BrokenReport
}

type BrokenReport struct {
	SourceKey       string
	ForeignKeyField string
	TargetTable     string
	TargetKeyField  string
	TargetKey       string
}

type fileState struct {
	index   int
	report  FileReport
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

func NewValidator(exportName, url string, textWriter, jsonWriter io.Writer) *Validator {
	return &Validator{
		textWriter: textWriter,
		jsonWriter: jsonWriter,
		key:        make(map[string]map[string]int),
		files:      make([]string, 0),
		errCount:   0,
		report: ExportReport{
			ExportName: exportName,
			Url:        url,
			Files:      make([]FileReport, 0),
		},
	}
}

// EnterTable implements Recorder.
func (v *Validator) EnterTable(td spec.Table) error {
	// Sanity check: there should be no table specific state.
	if v.td != nil {
		return fmt.Errorf("already processing a table, did you forget to call LeaveTable()?")
	}

	v.td = &td
	return nil
}

// LeaveTable implements Recorder.
func (v *Validator) LeaveTable() error {
	// Sanity check: there should be table specific state.
	if v.td == nil {
		return fmt.Errorf("not processing a table, did you forget to call EnterTable()?")
	}

	v.td = nil
	return nil
}

// EnterFile implements Recorder.
func (v *Validator) EnterFile(f string) error {
	// Sanity check: have we seen this file before?
	for _, name := range v.files {
		if name == f {
			return fmt.Errorf("file %s already validated", f)
		}
	}

	// Sanity check: there should be no file specific state.
	if v.fileState != nil {
		return fmt.Errorf("already processing a file, did you forget to call LeaveFile()?")
	}

	// Insert the file name into `v.files`.
	v.files = append(v.files, f)

	// Initialize file specific state.
	v.fileState = &fileState{
		report: FileReport{
			FileName: f,
			Broken:   make([]BrokenReport, 0),
		},
		index:   len(v.files) - 1,
	}

	fmt.Fprintln(v.textWriter, f)
	return nil
}

// LeaveFile implements Recorder.
func (v *Validator) LeaveFile() error {
	s := v.fileState

	// Sanity check: there should be file specific state.
	if s == nil {
		return fmt.Errorf("not processing a file, did you forget to call EnterFile()?")
	}

	// Append to the report.
	v.report.Files = append(v.report.Files, s.report)
	v.fileState = nil
	return nil
}

// Record implements Recorder.
func (v *Validator) Record(item map[string]string) error {
	td := v.td
	s := v.fileState
	fileName := s.report.FileName

	// Check for duplicate key definitions.
	key := item[td.Primary]
	if _, ok := v.key[td.Element]; !ok {
		v.key[td.Element] = make(map[string]int)
	}
	keys := v.key[td.Element]
	if originalFileIndex, ok := keys[key]; ok {
		v.reportDuplicate(duplicate{
			table:         td.Element,
			column:        td.Primary,
			key:           key,
			originalFile:  v.files[originalFileIndex],
			duplicateFile: fileName,
		})
		v.errCount++
	}
	keys[key] = s.index

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
			v.reportBroken(brk)
			v.errCount++

			s.report.Broken = append(s.report.Broken, BrokenReport{
				SourceKey:       key,
				ForeignKeyField: field.Name,
				TargetTable:     ref.Table,
				TargetKeyField:  ref.Column,
				TargetKey:       x,
			})
			s.report.NumBroken++

			continue
		}
		if _, ok := v.key[ref.Table][x]; !ok {
			v.reportBroken(brk)
			v.errCount++

			s.report.Broken = append(s.report.Broken, BrokenReport{
				SourceKey:       key,
				ForeignKeyField: field.Name,
				TargetTable:     ref.Table,
				TargetKeyField:  ref.Column,
				TargetKey:       x,
			})
			s.report.NumBroken++

			continue
		}
	}
	return nil
}

// Close implements io.Closer.
func (v *Validator) Close() error {
	if err := json.NewEncoder(v.jsonWriter).Encode(v.report); err != nil {
		return err
	}

	if v.errCount == 0 {
		fmt.Fprintln(v.textWriter, "SUCCESS")
	} else {
		fmt.Fprintf(v.textWriter, "FAILURE: %d error(s) found\n", v.errCount)
	}

	return nil
}

func (v *Validator) reportDuplicate(dup duplicate) {
	fmt.Fprintf(v.textWriter, "- duplicate: %s.%s=%s already appeared in %s\n", dup.table, dup.column, dup.key, dup.originalFile)
}

func (v *Validator) reportBroken(brk broken) {
	fmt.Fprintf(v.textWriter, "- broken: %s(%s=%s).%s references %s.%s=%s, which is missing\n", brk.table, brk.primary, brk.key, brk.column, brk.targetTable, brk.targetColumn, brk.targetKey)
}
