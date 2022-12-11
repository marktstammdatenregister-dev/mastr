package generator

import (
	"bytes"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func FormatAndWriteOutput(buf bytes.Buffer, outputPrefix string, outputSuffix string, fileName string) {
	src, err := format.Source(buf.Bytes())
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		src = buf.Bytes()
	}

	output := strings.ToLower(outputPrefix + fileName + outputSuffix + ".go")
	outputPath := filepath.Join(".", output)
	if err := os.WriteFile(outputPath, src, 0644); err != nil {
		log.Fatalf("writing output: %s", err)
	}
}
