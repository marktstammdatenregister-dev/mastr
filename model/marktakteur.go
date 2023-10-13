// Code generated by go mastr/model generator; DO NOT EDIT.

package main

import (
	"encoding/xml"
)

// Marktakteure
type Marktakteure struct {
	XMLName      xml.Name      `xml:"Marktakteure"`
	Marktakteure []Marktakteur `xml:"Marktakteur"`
}

type Marktakteur struct {
	XMLName                              xml.Name `xml:"Marktakteur"`
	MastrNummer                          string   `xml:"MastrNummer"`
	DatumLetzeAktualisierung             string   `xml:"DatumLetzeAktualisierung"`
	Personenart                          uint     `xml:"Personenart"` // referenziert Katalogwert
	MarktakteurVorname                   string   `xml:"MarktakteurVorname"`
	MarktakteurNachname                  string   `xml:"MarktakteurNachname"`
	Firmenname                           string   `xml:"Firmenname"`
	Marktfunktion                        uint     `xml:"Marktfunktion"`
	Rechtsform                           uint     `xml:"Rechtsform"` // referenziert Katalogwert
	SonstigeRechtsform                   string   `xml:"SonstigeRechtsform"`
	Marktrollen                          string   `xml:"Marktrollen"`
	Land                                 uint     `xml:"Land"` // referenziert Katalogwert
	Region                               string   `xml:"Region"`
	Strasse                              string   `xml:"Strasse"`
	Hausnummer                           string   `xml:"Hausnummer"`
	Hausnummer_nv                        string   `xml:"Hausnummer_nv"`
	Adresszusatz                         string   `xml:"Adresszusatz"`
	Postleitzahl                         string   `xml:"Postleitzahl"`
	Ort                                  string   `xml:"Ort"`
	Bundesland                           uint     `xml:"Bundesland"` // referenziert Katalogwert
	Netz                                 string   `xml:"Netz"`
	Nuts2                                string   `xml:"Nuts2"`
	Email                                string   `xml:"Email"`
	Telefon                              string   `xml:"Telefon"`
	Fax                                  string   `xml:"Fax"`
	Fax_nv                               string   `xml:"Fax_nv"`
	Webseite                             string   `xml:"Webseite"`
	Webseite_nv                          string   `xml:"Webseite_nv"`
	Registergericht                      uint     `xml:"Registergericht"`
	Registergericht_nv                   string   `xml:"Registergericht_nv"`
	RegistergerichtAusland               string   `xml:"RegistergerichtAusland"`
	RegistergerichtAusland_nv            string   `xml:"RegistergerichtAusland_nv"`
	Registernummer                       string   `xml:"Registernummer"`
	Registernummer_nv                    string   `xml:"Registernummer_nv"`
	RegisternummerAusland                string   `xml:"RegisternummerAusland"`
	RegisternummerAusland_nv             string   `xml:"RegisternummerAusland_nv"`
	Taetigkeitsbeginn                    string   `xml:"Taetigkeitsbeginn"`
	AcerCode                             string   `xml:"AcerCode"`
	AcerCode_nv                          string   `xml:"AcerCode_nv"`
	Umsatzsteueridentifikationsnummer    string   `xml:"Umsatzsteueridentifikationsnummer"`
	Umsatzsteueridentifikationsnummer_nv string   `xml:"Umsatzsteueridentifikationsnummer_nv"`
	Taetigkeitsende                      string   `xml:"Taetigkeitsende"`
	Taetigkeitsende_nv                   string   `xml:"Taetigkeitsende_nv"`
	BundesnetzagenturBetriebsnummer      string   `xml:"BundesnetzagenturBetriebsnummer"`
	BundesnetzagenturBetriebsnummer_nv   string   `xml:"BundesnetzagenturBetriebsnummer_nv"`
	LandAnZustelladresse                 uint     `xml:"LandAnZustelladresse"`
	PostleitzahlAnZustelladresse         string   `xml:"PostleitzahlAnZustelladresse"`
	OrtAnZustelladresse                  string   `xml:"OrtAnZustelladresse"`
	StrasseAnZustelladresse              string   `xml:"StrasseAnZustelladresse"`
	HausnummerAnZustelladresse           string   `xml:"HausnummerAnZustelladresse"`
	HausnummerAnZustelladresse_nv        string   `xml:"HausnummerAnZustelladresse_nv"`
	AdresszusatzAnZustelladresse         string   `xml:"AdresszusatzAnZustelladresse"`
	Kmu                                  string   `xml:"Kmu"`
	TelefonnummerVMav                    string   `xml:"TelefonnummerVMav"`
	EmailVMav                            string   `xml:"EmailVMav"`
	RegistrierungsdatumMarktakteur       string   `xml:"RegistrierungsdatumMarktakteur"`
	HauptwirtdschaftszweigAbteilung      uint     `xml:"HauptwirtdschaftszweigAbteilung"` // referenziert Katalogwert
	HauptwirtdschaftszweigGruppe         uint     `xml:"HauptwirtdschaftszweigGruppe"`    // referenziert Katalogwert
	HauptwirtdschaftszweigAbschnitt      uint     `xml:"HauptwirtdschaftszweigAbschnitt"` // referenziert Katalogwert
	Direktvermarktungsunternehmen        string   `xml:"Direktvermarktungsunternehmen"`
	BelieferungVonLetztverbrauchernStrom string   `xml:"BelieferungVonLetztverbrauchernStrom"`
	BelieferungHaushaltskundenStrom      string   `xml:"BelieferungHaushaltskundenStrom"`
	Gasgrosshaendler                     string   `xml:"Gasgrosshaendler"`
	Stromgrosshaendler                   string   `xml:"Stromgrosshaendler"`
	BelieferungVonLetztverbrauchernGas   string   `xml:"BelieferungVonLetztverbrauchernGas"`
	BelieferungHaushaltskundenGas        string   `xml:"BelieferungHaushaltskundenGas"`
}
