package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"os"
)

// ExampleReadXml liest eine XML aus dem Gesamtdownload des Marktstammdatenregisters
func ExampleReadXml() {
	xmlBytes, err := ReadFileUTF16("Katalogwerte.xml")
	if err != nil {
		panic(err)
	}
	var katalogwerte Katalogwerte
	err = xml.Unmarshal(xmlBytes, &katalogwerte)
	if err != nil {
		panic(err)
	}
	for _, katalogwert := range katalogwerte.Katalogwerte {
		println(fmt.Sprintf("%d: %s\n", katalogwert.KatalogKategorieId, katalogwert.Wert))
	}
}

// ReadFileUTF16 Similar to os.ReadFile() but decodes UTF-16.  Useful when
// reading data from MS-Windows systems that generate UTF-16BE files,
// but will do the right thing if other BOMs are found.
func ReadFileUTF16(filename string) ([]byte, error) {
	raw, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM)
	unicodeReader := transform.NewReader(bytes.NewReader(raw), utf16.NewDecoder())
	decoded, err := io.ReadAll(unicodeReader)
	return decoded, err
}
