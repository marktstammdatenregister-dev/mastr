// Code generated by go mastr/model generator; DO NOT EDIT.

package main

import (
	"encoding/xml"
)

// Katalogkategorien
type Katalogkategorien struct {
	XMLName           xml.Name           `xml:"Katalogkategorien"`
	Katalogkategorien []Katalogkategorie `xml:"Katalogkategorie"`
}

type Katalogkategorie struct {
	XMLName xml.Name `xml:"Katalogkategorie"`
	Id      uint     `xml:"Id"`
	Name    string   `xml:"Name"`
}
