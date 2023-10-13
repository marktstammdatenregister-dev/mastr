// Code generated by go mastr/model generator; DO NOT EDIT.

package main

import (
	"encoding/xml"
)

// AnlagenEegGeoSolarthermieGrubenKlaerschlammDruckentspannung
type AnlagenEegGeoSolarthermieGrubenKlaerschlammDruckentspannung struct {
	XMLName                                                     xml.Name                                                     `xml:"AnlagenEegGeoSolarthermieGrubenKlaerschlammDruckentspannung"`
	AnlagenEegGeoSolarthermieGrubenKlaerschlammDruckentspannung []AnlageEegGeoSolarthermieGrubenKlaerschlammDruckentspannung `xml:"AnlageEegGeoSolarthermieGrubenKlaerschlammDruckentspannung"`
}

type AnlageEegGeoSolarthermieGrubenKlaerschlammDruckentspannung struct {
	XMLName                             xml.Name `xml:"AnlageEegGeoSolarthermieGrubenKlaerschlammDruckentspannung"`
	Registrierungsdatum                 string   `xml:"Registrierungsdatum"`
	DatumLetzteAktualisierung           string   `xml:"DatumLetzteAktualisierung"`
	EegInbetriebnahmedatum              string   `xml:"EegInbetriebnahmedatum"`
	EegMaStRNummer                      string   `xml:"EegMaStRNummer"`
	AnlagenschluesselEeg                string   `xml:"AnlagenschluesselEeg"`
	AnlagenkennzifferAnlagenregister    string   `xml:"AnlagenkennzifferAnlagenregister"`
	AnlagenkennzifferAnlagenregister_nv string   `xml:"AnlagenkennzifferAnlagenregister_nv"`
	InstallierteLeistung                float32  `xml:"InstallierteLeistung"`
	AnlageBetriebsstatus                uint     `xml:"AnlageBetriebsstatus"` // referenziert Katalogwert
	VerknuepfteEinheitenMaStRNummern    string   `xml:"VerknuepfteEinheitenMaStRNummern"`
}
