// Code generated by go mastr/model generator; DO NOT EDIT.

package main

import (
	"encoding/xml"
)

// EinheitenVerbrennung
type EinheitenVerbrennung struct {
	XMLName              xml.Name             `xml:"EinheitenVerbrennung"`
	EinheitenVerbrennung []EinheitVerbrennung `xml:"EinheitVerbrennung"`
}

type EinheitVerbrennung struct {
	XMLName                                            xml.Name `xml:"EinheitVerbrennung"`
	EinheitMastrNummer                                 string   `xml:"EinheitMastrNummer"`
	DatumLetzteAktualisierung                          string   `xml:"DatumLetzteAktualisierung"`
	LokationMaStRNummer                                string   `xml:"LokationMaStRNummer"` // referenziert Lokation
	NetzbetreiberpruefungStatus                        uint     `xml:"NetzbetreiberpruefungStatus"`
	NetzbetreiberpruefungDatum                         string   `xml:"NetzbetreiberpruefungDatum"`
	AnlagenbetreiberMastrNummer                        string   `xml:"AnlagenbetreiberMastrNummer"` // referenziert Marktakteur
	Land                                               uint     `xml:"Land"`                        // referenziert Katalogwert
	Bundesland                                         uint     `xml:"Bundesland"`                  // referenziert Katalogwert
	Landkreis                                          string   `xml:"Landkreis"`
	Gemeinde                                           string   `xml:"Gemeinde"`
	Gemeindeschluessel                                 string   `xml:"Gemeindeschluessel"`
	Postleitzahl                                       string   `xml:"Postleitzahl"`
	Gemarkung                                          string   `xml:"Gemarkung"`
	FlurFlurstuecknummern                              string   `xml:"FlurFlurstuecknummern"`
	Strasse                                            string   `xml:"Strasse"`
	StrasseNichtGefunden                               string   `xml:"StrasseNichtGefunden"`
	Hausnummer                                         string   `xml:"Hausnummer"`
	Hausnummer_nv                                      string   `xml:"Hausnummer_nv"`
	HausnummerNichtGefunden                            string   `xml:"HausnummerNichtGefunden"`
	Adresszusatz                                       string   `xml:"Adresszusatz"`
	Ort                                                string   `xml:"Ort"`
	Laengengrad                                        float32  `xml:"Laengengrad"`
	Breitengrad                                        float32  `xml:"Breitengrad"`
	Registrierungsdatum                                string   `xml:"Registrierungsdatum"`
	GeplantesInbetriebnahmedatum                       string   `xml:"GeplantesInbetriebnahmedatum"`
	Inbetriebnahmedatum                                string   `xml:"Inbetriebnahmedatum"`
	DatumEndgueltigeStilllegung                        string   `xml:"DatumEndgueltigeStilllegung"`
	DatumBeginnVoruebergehendeStilllegung              string   `xml:"DatumBeginnVoruebergehendeStilllegung"`
	DatumWiederaufnahmeBetrieb                         string   `xml:"DatumWiederaufnahmeBetrieb"`
	EinheitSystemstatus                                uint     `xml:"EinheitSystemstatus"`   // referenziert Katalogwert
	EinheitBetriebsstatus                              uint     `xml:"EinheitBetriebsstatus"` // referenziert Katalogwert
	BestandsanlageMastrNummer                          string   `xml:"BestandsanlageMastrNummer"`
	NichtVorhandenInMigriertenEinheiten                string   `xml:"NichtVorhandenInMigriertenEinheiten"`
	AltAnlagenbentreiberMastrNummer                    string   `xml:"AltAnlagenbentreiberMastrNummer"`
	DatumDesBetreiberwechsels                          string   `xml:"DatumDesBetreiberwechsels"`
	DatumRegistrierungDesBetreiberwechsels             string   `xml:"DatumRegistrierungDesBetreiberwechsels"`
	NameStromerzeugungseinheit                         string   `xml:"NameStromerzeugungseinheit"`
	Weic                                               string   `xml:"Weic"`
	Weic_nv                                            string   `xml:"Weic_nv"`
	WeicDisplayName                                    string   `xml:"WeicDisplayName"`
	Kraftwerksnummer                                   string   `xml:"Kraftwerksnummer"`
	Kraftwerksnummer_nv                                string   `xml:"Kraftwerksnummer_nv"`
	Energietraeger                                     uint     `xml:"Energietraeger"` // referenziert Katalogwert
	Bruttoleistung                                     float32  `xml:"Bruttoleistung"`
	Nettonennleistung                                  float32  `xml:"Nettonennleistung"`
	AnschlussAnHoechstOderHochSpannung                 string   `xml:"AnschlussAnHoechstOderHochSpannung"`
	Schwarzstartfaehigkeit                             string   `xml:"Schwarzstartfaehigkeit"`
	Inselbetriebsfaehigkeit                            string   `xml:"Inselbetriebsfaehigkeit"`
	Einsatzverantwortlicher                            string   `xml:"Einsatzverantwortlicher"`
	FernsteuerbarkeitNb                                string   `xml:"FernsteuerbarkeitNb"`
	FernsteuerbarkeitDv                                string   `xml:"FernsteuerbarkeitDv"`
	FernsteuerbarkeitDr                                string   `xml:"FernsteuerbarkeitDr"`
	Einspeisungsart                                    uint     `xml:"Einspeisungsart"` // referenziert Katalogwert
	PraequalifiziertFuerRegelenergie                   string   `xml:"PraequalifiziertFuerRegelenergie"`
	GenMastrNummer                                     string   `xml:"GenMastrNummer"` // referenziert EinheitGenehmigung
	NameKraftwerk                                      string   `xml:"NameKraftwerk"`
	NameKraftwerksblock                                string   `xml:"NameKraftwerksblock"`
	DatumBaubeginn                                     string   `xml:"DatumBaubeginn"`
	AnzeigeEinerStilllegung                            string   `xml:"AnzeigeEinerStilllegung"`
	ArtDerStilllegung                                  uint     `xml:"ArtDerStilllegung"` // referenziert Katalogwert
	DatumBeginnVorlaeufigenOderEndgueltigenStilllegung string   `xml:"DatumBeginnVorlaeufigenOderEndgueltigenStilllegung"`
	SteigerungNettonennleistungKombibetrieb            float32  `xml:"SteigerungNettonennleistungKombibetrieb"`
	AnlageIstImKombibetrieb                            string   `xml:"AnlageIstImKombibetrieb"`
	MastrNummernKombibetrieb                           string   `xml:"MastrNummernKombibetrieb"`
	NetzreserveAbDatum                                 string   `xml:"NetzreserveAbDatum"`
	SicherheitsbereitschaftAbDatum                     string   `xml:"SicherheitsbereitschaftAbDatum"`
	Hauptbrennstoff                                    uint     `xml:"Hauptbrennstoff"`         // referenziert Katalogwert
	WeitererHauptbrennstoff                            uint     `xml:"WeitererHauptbrennstoff"` // referenziert Katalogwert
	WeitereBrennstoffe                                 string   `xml:"WeitereBrennstoffe"`
	BestandteilGrenzkraftwerk                          string   `xml:"BestandteilGrenzkraftwerk"`
	NettonennleistungDeutschland                       float32  `xml:"NettonennleistungDeutschland"`
	AnteiligNutzungsberechtigte                        string   `xml:"AnteiligNutzungsberechtigte"`
	Notstromaggregat                                   string   `xml:"Notstromaggregat"`
	Einsatzort                                         uint     `xml:"Einsatzort"`     // referenziert Katalogwert
	KwkMaStRNummer                                     string   `xml:"KwkMaStRNummer"` // referenziert AnlageKwk
	Technologie                                        uint     `xml:"Technologie"`    // referenziert Katalogwert

}
