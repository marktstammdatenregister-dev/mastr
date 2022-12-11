// Code generated by go mastr/model generator; DO NOT EDIT.

package main

import (
	"encoding/xml"
)

// EinheitenStromSpeicher
type EinheitenStromSpeicher struct {
	XMLName                xml.Name               `xml:"EinheitenStromSpeicher"`
	EinheitenStromSpeicher []EinheitStromSpeicher `xml:"EinheitStromSpeicher"`
}

type EinheitStromSpeicher struct {
	XMLName                                xml.Name `xml:"EinheitStromSpeicher"`
	EinheitMastrNummer                     string   `xml:"EinheitMastrNummer"`
	DatumLetzteAktualisierung              string   `xml:"DatumLetzteAktualisierung"`
	LokationMaStRNummer                    string   `xml:"LokationMaStRNummer"` // referenziert Lokation
	NetzbetreiberpruefungStatus            uint     `xml:"NetzbetreiberpruefungStatus"`
	NetzbetreiberpruefungDatum             string   `xml:"NetzbetreiberpruefungDatum"`
	AnlagenbetreiberMastrNummer            string   `xml:"AnlagenbetreiberMastrNummer"` // referenziert Marktakteur
	Land                                   uint     `xml:"Land"`                        // referenziert Katalogwert
	Bundesland                             uint     `xml:"Bundesland"`                  // referenziert Katalogwert
	Landkreis                              string   `xml:"Landkreis"`
	Gemeinde                               string   `xml:"Gemeinde"`
	Gemeindeschluessel                     string   `xml:"Gemeindeschluessel"`
	Postleitzahl                           string   `xml:"Postleitzahl"`
	Gemarkung                              string   `xml:"Gemarkung"`
	FlurFlurstuecknummern                  string   `xml:"FlurFlurstuecknummern"`
	Strasse                                string   `xml:"Strasse"`
	StrasseNichtGefunden                   string   `xml:"StrasseNichtGefunden"`
	Hausnummer                             string   `xml:"Hausnummer"`
	Hausnummer_nv                          string   `xml:"Hausnummer_nv"`
	HausnummerNichtGefunden                string   `xml:"HausnummerNichtGefunden"`
	Adresszusatz                           string   `xml:"Adresszusatz"`
	Ort                                    string   `xml:"Ort"`
	Laengengrad                            float32  `xml:"Laengengrad"`
	Breitengrad                            float32  `xml:"Breitengrad"`
	Registrierungsdatum                    string   `xml:"Registrierungsdatum"`
	GeplantesInbetriebnahmedatum           string   `xml:"GeplantesInbetriebnahmedatum"`
	Inbetriebnahmedatum                    string   `xml:"Inbetriebnahmedatum"`
	DatumEndgueltigeStilllegung            string   `xml:"DatumEndgueltigeStilllegung"`
	DatumBeginnVoruebergehendeStilllegung  string   `xml:"DatumBeginnVoruebergehendeStilllegung"`
	DatumWiederaufnahmeBetrieb             string   `xml:"DatumWiederaufnahmeBetrieb"`
	EinheitSystemstatus                    uint     `xml:"EinheitSystemstatus"`   // referenziert Katalogwert
	EinheitBetriebsstatus                  uint     `xml:"EinheitBetriebsstatus"` // referenziert Katalogwert
	BestandsanlageMastrNummer              string   `xml:"BestandsanlageMastrNummer"`
	NichtVorhandenInMigriertenEinheiten    string   `xml:"NichtVorhandenInMigriertenEinheiten"`
	AltAnlagenbentreiberMastrNummer        string   `xml:"AltAnlagenbentreiberMastrNummer"`
	DatumDesBetreiberwechsels              string   `xml:"DatumDesBetreiberwechsels"`
	DatumRegistrierungDesBetreiberwechsels string   `xml:"DatumRegistrierungDesBetreiberwechsels"`
	NameStromerzeugungseinheit             string   `xml:"NameStromerzeugungseinheit"`
	Weic                                   string   `xml:"Weic"`
	Weic_nv                                string   `xml:"Weic_nv"`
	WeicDisplayName                        string   `xml:"WeicDisplayName"`
	Kraftwerksnummer                       string   `xml:"Kraftwerksnummer"`
	Kraftwerksnummer_nv                    string   `xml:"Kraftwerksnummer_nv"`
	Energietraeger                         uint     `xml:"Energietraeger"` // referenziert Katalogwert
	Bruttoleistung                         float32  `xml:"Bruttoleistung"`
	Nettonennleistung                      float32  `xml:"Nettonennleistung"`
	AnschlussAnHoechstOderHochSpannung     string   `xml:"AnschlussAnHoechstOderHochSpannung"`
	Schwarzstartfaehigkeit                 string   `xml:"Schwarzstartfaehigkeit"`
	Inselbetriebsfaehigkeit                string   `xml:"Inselbetriebsfaehigkeit"`
	Einsatzverantwortlicher                string   `xml:"Einsatzverantwortlicher"`
	FernsteuerbarkeitNb                    string   `xml:"FernsteuerbarkeitNb"`
	FernsteuerbarkeitDv                    string   `xml:"FernsteuerbarkeitDv"`
	FernsteuerbarkeitDr                    string   `xml:"FernsteuerbarkeitDr"`
	Einspeisungsart                        uint     `xml:"Einspeisungsart"` // referenziert Katalogwert
	PraequalifiziertFuerRegelenergie       string   `xml:"PraequalifiziertFuerRegelenergie"`
	GenMastrNummer                         string   `xml:"GenMastrNummer"`      // referenziert EinheitGenehmigung
	Einsatzort                             uint     `xml:"Einsatzort"`          // referenziert Katalogwert
	AcDcKoppelung                          uint     `xml:"AcDcKoppelung"`       // referenziert Katalogwert
	Batterietechnologie                    uint     `xml:"Batterietechnologie"` // referenziert Katalogwert
	PumpbetriebLeistungsaufnahme           float32  `xml:"PumpbetriebLeistungsaufnahme"`
	PumpbetriebKontinuierlichRegelbar      string   `xml:"PumpbetriebKontinuierlichRegelbar"`
	Pumpspeichertechnologie                uint     `xml:"Pumpspeichertechnologie"` // referenziert Katalogwert
	Notstromaggregat                       string   `xml:"Notstromaggregat"`
	BestandteilGrenzkraftwerk              string   `xml:"BestandteilGrenzkraftwerk"`
	NettonennleistungDeutschland           float32  `xml:"NettonennleistungDeutschland"`
	ZugeordnenteWirkleistungWechselrichter float32  `xml:"ZugeordnenteWirkleistungWechselrichter"`
	SpeMastrNummer                         string   `xml:"SpeMastrNummer"` // referenziert AnlageStromSpeicher
	EegMaStRNummer                         string   `xml:"EegMaStRNummer"`
	EegAnlagentyp                          uint     `xml:"EegAnlagentyp"` // referenziert Einheitentyp
	Technologie                            uint     `xml:"Technologie"`   // referenziert Katalogwert

}
