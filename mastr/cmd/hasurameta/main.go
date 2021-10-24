package main

import (
	"flag"
	"log"
	"pvdb.de/mastr/internal"
)

func main() {
	const defaultOption = "<undefined>"
	descriptorFileName := flag.String("descriptor", defaultOption, "file name of the table descriptor")
	schemaName := flag.String("schema", defaultOption, "schema name")
	flag.Parse()

	td, err := internal.DecodeDescriptor(*descriptorFileName)
	if err != nil {
		log.Fatalf("failed to decode descriptor: %w", err)
	}

	b, err := internal.ToHasura(*schemaName, *td)
	if err != nil {
		log.Fatalf("failed to encode metadata: %w", err)
	}
	log.Printf("%s", string(b))
}
