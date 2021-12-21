package main

import (
	"flag"
	"log"
	"pvdb.de/mastr/internal"
	"pvdb.de/mastr/internal/spec"
)

func main() {
	const defaultOption = "<undefined>"
	specFileName := flag.String("spec", defaultOption, "file name of the table spec")
	schemaName := flag.String("schema", defaultOption, "schema name")
	flag.Parse()

	td, err := spec.DecodeTable(*specFileName)
	if err != nil {
		log.Fatalf("failed to decode spec: %w", err)
	}

	b, err := internal.ToHasura(*schemaName, *td)
	if err != nil {
		log.Fatalf("failed to encode metadata: %w", err)
	}
	log.Printf("%s", string(b))
}
