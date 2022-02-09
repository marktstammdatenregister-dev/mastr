package internal

import (
	"marktstammdatenregister.dev/internal/spec"
)

type Recorder interface {
	EnterTable(spec.Table) error
	LeaveTable() error
	EnterFile(string) error
	LeaveFile() error
	Record(map[string]string) error
}
