CREATE TABLE "Marktakteure" (
  'MastrNummer' TEXT PRIMARY KEY,
  'DatumLetzeAktualisierung' TEXT,
  'Personenart' TEXT NOT NULL,
  'MarktakteurVorname' TEXT,
  'MarktakteurNachname' TEXT,
  'Firmenname' TEXT,
  'Marktfunktion' INTEGER NOT NULL,
  'Rechtsform' INTEGER,
  'SonstigeRechtsform' TEXT,
  'Marktrollen' TEXT,
  'Land' INTEGER,
  'Region' TEXT,
  'Strasse' TEXT,
  'Hausnummer' TEXT,
  'Hausnummer_nv' INTEGER,
  'Adresszusatz' TEXT,
  'Postleitzahl' TEXT,
  'Ort' TEXT,
  'Bundesland' INTEGER,
  'Netz' TEXT,
  'Nuts2' TEXT,
  'Email' TEXT,
  'Telefon' TEXT,
  'Fax' TEXT,
  'Fax_nv' INTEGER,
  'Webseite' TEXT,
  'Webseite_nv' INTEGER,
  'Registergericht' INTEGER,
  'Registergericht_nv' INTEGER,
  'RegistergerichtAusland' TEXT,
  'RegistergerichtAusland_nv' INTEGER,
  'Registernummer' TEXT,
  'Registernummer_nv' INTEGER,
  'RegisternummerAusland' TEXT,
  'RegisternummerAusland_nv' INTEGER,
  'Taetigkeitsbeginn' TEXT,
  'AcerCode' TEXT,
  'AcerCode_nv' INTEGER,
  'Umsatzsteueridentifikationsnummer' TEXT,
  'Umsatzsteueridentifikationsnummer_nv' INTEGER,
  'Taetigkeitsende' TEXT,
  'Taetigkeitsende_nv' INTEGER,
  'BundesnetzagenturBetriebsnummer' TEXT,
  'BundesnetzagenturBetriebsnummer_nv' INTEGER
  'LandAnZustelladresse' INTEGER,
  'PostleitzahlAnZustelladresse' TEXT,
  'OrtAnZustelladresse' TEXT,
  'StrasseAnZustelladresse' TEXT,
  'HausnummerAnZustelladresse' TEXT,
  'HausnummerAnZustelladresse_nv' INTEGER,
  'AdresszusatzAnZustelladresse' TEXT,
  'Kmu' INTEGER,
  'TelefonnummerVMav' TEXT,
  'EmailVMav' TEXT,
  'RegistrierungsdatumMarktakteur' TEXT,
  'HauptwirtdschaftszweigAbteilung' INTEGER,
  'HauptwirtdschaftszweigGruppe' INTEGER,
  'HauptwirtdschaftszweigAbschnitt' INTEGER,
  'Direktvermarktungsunternehmen' INTEGER,
  'BelieferungVonLetztverbrauchernStrom' INTEGER,
  'BelieferungHaushaltskundenStrom' INTEGER,
  'Gasgrosshaendler' INTEGER,
  'Stromgrosshaendler' INTEGER,
  'BelieferungVonLetztverbrauchernGas' INTEGER,
  'BelieferungHaushaltskundenGas' INTEGER
);
