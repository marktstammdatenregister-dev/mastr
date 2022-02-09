package internal

import (
	"marktstammdatenregister.dev/internal/spec"
)

type Recorder interface {
	EnterTable(spec.Table)
	LeaveTable()
	EnterFile(string)
	LeaveFile()
	Record(map[string]string) error
}
