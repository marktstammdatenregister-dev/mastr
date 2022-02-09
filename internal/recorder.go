package internal

import (
	"io"
	"marktstammdatenregister.dev/internal/spec"
)

type Recorder interface {
	io.Closer
	EnterTable(spec.Table) error
	LeaveTable() error
	EnterFile(string) error
	LeaveFile() error
	Record(map[string]string) error
}
