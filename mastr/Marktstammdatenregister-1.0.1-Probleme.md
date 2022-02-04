# Probleme mit dem Marktstammdatenregister Gesamtdatenexport

Das [Marktstammdatenregister][mastr] ist das Register für den deutschen Strom- und Gasmarkt.
Es wird von der Bundesnetzagentur geführt.
Alle Anlagebetreiber, Netzbetreiber, Energielieferanten sowie alle Strom- und Gaserzeugungsanlagen sind darin registriert.
Das macht es zu einer enorm wichtigen Datenquelle, um die Entwicklung der deutschen Energieversorgungskapazitäten zu verstehen.

Auf der Website des Marktstammdatenregisters gibt es eine Filterfunktion, deren Funktionsumfang aber eingeschränkt ist.
Von den Suchergebnissen lassen sich immer nur die ersten 5000 Datensätze herunterladen.

Als Alternative gibt es den Gesamtdatenexport [zum Download][export-download].
Zumindest theoretisch ermöglichen die darin enthaltenen Rohdaten eine flexible Verarbeitung der Registereinträge.

# Struktur der Gesamtdatenexport-Datei

Der Gesamtdatenexport ist eine ZIP-Datei, die wiederum eine oder mehrere XML-Dateien pro Datentyp enthält.
Als Datentyp zählen hier zum Beispiel Windenergieeinheiten oder Biomasseanlagen.

Die [Dokumentation][export-doc-live] des Gesamtdatenexports ist auf der Website des Marktstammdatenregisters verfügbar.
Für dieses Dokument ist Version 1.0.1 vom 22. November 2021 die Grundlage ([hier archiviert][export-doc-archived]).

Um ein konkretes Beispiel zu nehmen: der Gesamtdatenexport vom 3. Februar 2022 ([hier archiviert][export-archived]) ist 931 MiB gross und enthält 146 XML-Dateien, die entpackt insgesamt 18,2 GiB in Anspruch nehmen.

```console
$ unzip -l Gesamtdatenexport_20220203__840cfde7b693453982d28db827025ff0.zip
Archive:  /home/leo/Downloads/Gesamtdatenexport_20220203__840cfde7b693453982d28db827025ff0.zip
  Length      Date    Time    Name
---------  ---------- -----   ----
 29103018  02-03-2022 01:00   AnlagenEegBiomasse.xml
   251906  02-03-2022 01:00   AnlagenEegGeoSolarthermieGrubenKlaerschlammDruckentspannung.xml
157154966  02-03-2022 01:00   AnlagenEegSolar_1.xml [...]
 91854182  02-03-2022 01:21   AnlagenEegSolar_23.xml
 73000088  02-03-2022 01:21   AnlagenEegSpeicher_1.xml [...]
 21445298  02-03-2022 01:21   AnlagenEegSpeicher_4.xml
  6979416  02-03-2022 01:21   AnlagenEegWasser.xml
 68619008  02-03-2022 01:21   AnlagenEegWind.xml
    41106  02-03-2022 02:05   AnlagenGasSpeicher.xml
 65723506  02-03-2022 02:09   AnlagenKwk.xml
 82669008  02-03-2022 02:04   AnlagenStromSpeicher_1.xml [...]
 33628818  02-03-2022 02:05   AnlagenStromSpeicher_4.xml
   524404  02-03-2022 02:05   Bilanzierungsgebiete.xml
 79146124  02-03-2022 01:22   EinheitenBiomasse.xml
   804182  02-03-2022 01:47   EinheitenGasErzeuger.xml
   175534  02-03-2022 01:47   EinheitenGasSpeicher.xml
  1994400  02-03-2022 01:47   EinheitenGasverbraucher.xml
 30562806  02-03-2022 02:05   EinheitenGenehmigung.xml
  1601668  02-03-2022 01:44   EinheitenGeoSolarthermieGrubenKlaerschlammDruckentspannung.xml
    21108  02-03-2022 01:44   EinheitenKernkraft.xml
364861586  02-03-2022 01:22   EinheitenSolar_1.xml [...]
244848700  02-03-2022 01:44   EinheitenSolar_23.xml
326517904  02-03-2022 01:45   EinheitenStromSpeicher_1.xml [...]
133246692  02-03-2022 01:47   EinheitenStromSpeicher_4.xml
   918284  02-03-2022 01:47   EinheitenStromVerbraucher.xml
     1814  02-03-2022 01:47   Einheitentypen.xml
244167754  02-03-2022 01:44   EinheitenVerbrennung.xml
 28602896  02-03-2022 01:44   EinheitenWasser.xml
145815632  02-03-2022 01:22   EinheitenWind.xml
    29092  02-03-2022 02:04   Katalogkategorien.xml
   694514  02-03-2022 02:05   Katalogwerte.xml
 71403252  02-03-2022 02:06   Lokationen_1.xml [...]
 14905978  02-03-2022 02:09   Lokationen_25.xml
 64565014  02-03-2022 01:52   Marktakteure_1.xml [...]
 27554026  02-03-2022 02:04   Marktakteure_20.xml
  8508028  02-03-2022 02:05   Marktrollen.xml
120703968  02-03-2022 01:47   Netzanschlusspunkte_1.xml [...]
 94790298  02-03-2022 01:52   Netzanschlusspunkte_20.xml
   988162  02-03-2022 02:05   Netze.xml
---------                     -------
19535539596                     146 files
```

# Probleme der internen Konsistenz

Die Einträge im Gesamtdatenexport sind auf zwei verschiedene Arten intern nicht konsistent:

- Einträge sind dupliziert.
- Es wird auf Einträge verwiesen, die im Gesamtdatenexport nicht vorkommen.

## Duplizierte Einträge

Duplizierte Einträge kommen am Ende und am Anfang von aufeinander folgenden XML-Dateien eines bestimmten Datentyps vor.
Im Gesamtdatenexport vom 3. Februar 2022 findet sich ein Beispiel: die EEG-Solaranlage mit der Marktstammdatenregisternummer EEG930980053027 kommt sowohl als letzter Eintrag in `AnlagenEegSolar_7.xml` als auch als erster Eintrag in `AnlagenEegSolar_8` vor.

Um das nachzuvollziehen, formatieren wir die XML-Dateien zuerst leserlich mit `xmllint -format`.
Dann fügt `nl` Zeilennummern hinzu.
Als letztes suchen wir mit `grep` nach der Marktstammdatenregisternummer.

```
# Die leserlich formatierte Version von AnlagenEegSolar_7.xml hat 1378525 Zeilen.
$ xmllint -format AnlagenEegSolar_7.xml | wc --lines
1378525

# Die Anlage mit Marktstammdatenregisternummer EEG930980053027 ist der letzte Eintrag in AnlagenEegSolar_7.xml.
$ xmllint -format AnlagenEegSolar_7.xml | nl | grep --before-context=4 --after-context=9 EEG930980053027
1378512   <AnlageEegSolar>
1378513     <Registrierungsdatum>2021-01-02</Registrierungsdatum>
1378514     <DatumLetzteAktualisierung>2021-01-02T11:30:05.1746875</DatumLetzteAktualisierung>
1378515     <EegInbetriebnahmedatum>2011-10-25</EegInbetriebnahmedatum>
1378516     <EegMaStRNummer>EEG930980053027</EegMaStRNummer>
1378517     <AnlagenschluesselEeg>E2056801S160000000000022676200001</AnlagenschluesselEeg>
1378518     <AnlagenkennzifferAnlagenregister_nv>0</AnlagenkennzifferAnlagenregister_nv>
1378519     <InstallierteLeistung>25.480</InstallierteLeistung>
1378520     <RegistrierungsnummerPvMeldeportal_nv>0</RegistrierungsnummerPvMeldeportal_nv>
1378521     <AusschreibungZuschlag>0</AusschreibungZuschlag>
1378522     <AnlageBetriebsstatus>35</AnlageBetriebsstatus>
1378523     <VerknuepfteEinheitenMaStRNummern>SEE950922677951</VerknuepfteEinheitenMaStRNummern>
1378524   </AnlageEegSolar>
1378525 </AnlagenEegSolar>

# Die Anlage mit Marktstammdatenregisternummer EEG930980053027 ist der erste Eintrag in AnlagenEegSolar_8.xml.
$ xmllint -format AnlagenEegSolar_8.xml | nl | grep --before-context=6 --after-context=8 EEG930980053027
     1  <?xml version="1.0"?>
     2  <AnlagenEegSolar>
     3    <AnlageEegSolar>
     4      <Registrierungsdatum>2021-01-02</Registrierungsdatum>
     5      <DatumLetzteAktualisierung>2021-01-02T11:30:05.1746875</DatumLetzteAktualisierung>
     6      <EegInbetriebnahmedatum>2011-10-25</EegInbetriebnahmedatum>
     7      <EegMaStRNummer>EEG930980053027</EegMaStRNummer>
     8      <AnlagenschluesselEeg>E2056801S160000000000022676200001</AnlagenschluesselEeg>
     9      <AnlagenkennzifferAnlagenregister_nv>0</AnlagenkennzifferAnlagenregister_nv>
    10      <InstallierteLeistung>25.480</InstallierteLeistung>
    11      <RegistrierungsnummerPvMeldeportal_nv>0</RegistrierungsnummerPvMeldeportal_nv>
    12      <AusschreibungZuschlag>0</AusschreibungZuschlag>
    13      <AnlageBetriebsstatus>35</AnlageBetriebsstatus>
    14      <VerknuepfteEinheitenMaStRNummern>SEE950922677951</VerknuepfteEinheitenMaStRNummern>
    15    </AnlageEegSolar>
```

Eine Verdoppelung von Daten dieser Art sollte nicht passieren und deutet darauf hin, dass beim Exportprozess etwas schief läuft.

## Kaputte Verweise

Im Gesamtdatenexport vom 3. Februar 2022 gibt es mehrere Beispiele von Verweisen auf Einträge, die nicht im Export vorkommen.
Hier zwei Beispiele.

Erstes Beispiel: die Solareinheit mit Marktstammdatenregisternummer SEE956940074659 gibt als Anlagenbetreiber den Marktakteur mit Marktstammdatenregisternummer ABR930553771594 an, doch dieser findet sich nicht im Export.

```
# Marktstammdatenregisternummer ABR930553771594 kommt nur als Verweis vor, der Eintrag selbst fehlt.
$ grep -l ABR930553771594 *.xml
EinheitenSolar_13.xml

# Der Eintrag der verweisenden Solareinheit (gekürzt).
$ xmllint -format EinheitenSolar_13.xml | grep --before-context=6 --after-context=31 ABR930553771594
  <EinheitSolar>
    <EinheitMastrNummer>SEE956940074659</EinheitMastrNummer>
    <DatumLetzteAktualisierung>2021-12-17T10:19:00.3860556</DatumLetzteAktualisierung>
    <!-- ... -->
    <AnlagenbetreiberMastrNummer>ABR930553771594</AnlagenbetreiberMastrNummer>
    <!-- ... -->
  </EinheitSolar>
```

Unter "Anlagenbetreiber der Einheit" steht auf der Website des Marktstammdatenregisters im [Eintrag][broken-ex1] (archiviert am 3. Februar 2022 um 22:25 Uhr) der Solareinheit "deaktivierter Marktakteur (ABR930553771594)".
Das weist darauf hin, dass der Eintrag für den Marktakteur absichtlich nicht im Gesamtdatenexport vorkommt.
Meiner Meinung nach ist das Löschen oder Auslassen auf diese Art und Weise keine gute Methode.
Besser wäre es, Platzhaltereinträge für deaktivierte Marktakteure im Gesamtdatenexport zu erfassen.
Damit liessen sich deaktivierte und fehlende Marktakteure klar unterscheiden.

Zweites Beispiel: die Solareinheit mit Marktstammdatenregisternummer SEE983580197295 gibt EEG914024915342 als Marktstammdatenregisternummer der zugeordneten EEG-Anlage an, doch diese findet sich nicht im Export.

```
# Marktstammdatenregisternummer EEG914024915342 kommt nur als Verweis vor, der Eintrag selbst fehlt.
$ grep -l EEG914024915342 *.xml
EinheitenSolar_19.xml

# Der Eintrag der verweisenden Solareinheit (gekürzt).
$ xmllint -format EinheitenSolar_19.xml | grep --before-context=38 --after-context=1 EEG914024915342
  <EinheitSolar>
    <EinheitMastrNummer>SEE983580197295</EinheitMastrNummer>
    <DatumLetzteAktualisierung>2022-02-03T00:03:48.5877980</DatumLetzteAktualisierung>
    <!-- ... -->
    <EegMaStRNummer>EEG914024915342</EegMaStRNummer>
  </EinheitSolar>
```

Auf der Website des Marktstammdatenregisters [findet man][broken-ex2] (archiviert am 3. Februar 2022 um 22:26 Uhr) unter dem Reiter "EEG-Anlage" der Solareinheit Informationen zur Anlage EEG914024915342 (Stand 3. Februar 2022, 20:00 Uhr).
Interessant ist, dass die letzte Aktualisierung um 00:03 am 3. Februar stand, also 57 Minuten vor Beginn des Exportprozesses (angenommen, die Zeitstempel der XML-Dateien stimmen).
Dies könnte ein Hinweis darauf sein, dass der Exportprozess nicht durchweg den aktuellsten Stand der Daten erfasst.

# Unstimmigkeiten zwischen der Dokumentation und den Daten

## Fehlende Pflichtfelder

Die Dokumentation, die das Format der Gesamtdatenexporte beschreibt, definiert für jedes Feld, ob es ein Pflichtfeld ist.
Im Gesamtdatenexport fehlen viele dieser Pflichtfelder.

Die folgenden Pflichtfelder fehlem im Gesamtdatenexport vom 3. Februar 2022 (pro Datentyp ist immer nur eine XML-Datei angegeben):

```
Bilanzierungsgebiete.xml
- `BilanzierungsgebietNetzanschlusspunkt` fehlt in 498 Einträgen, z.B. Id=2
- `Yeic` fehlt in 3 Einträgen, z.B. Id=826
```
```
Marktakteure_1.xml
- `AcerCode` fehlt in 99904 Einträgen, z.B. MastrNummer=ABR111222333123
- `AcerCode_nv` fehlt in 90662 Einträgen, z.B. MastrNummer=ABR111222333123
- `BundesnetzagenturBetriebsnummer` fehlt in 89859 Einträgen, z.B. MastrNummer=ABR111222333123
- `BundesnetzagenturBetriebsnummer_nv` fehlt in 89859 Einträgen, z.B. MastrNummer=ABR111222333123
- `DatumLetzeAktualisierung` fehlt in 89859 Einträgen, z.B. MastrNummer=ABR111222333123
- `Email` fehlt in 92433 Einträgen, z.B. MastrNummer=ABR111222333123
- `Fax` fehlt in 97185 Einträgen, z.B. MastrNummer=ABR111222333123
- `Fax_nv` fehlt in 90662 Einträgen, z.B. MastrNummer=ABR111222333123
- `Hausnummer` fehlt in 89994 Einträgen, z.B. MastrNummer=ABR111222333123
- `Hausnummer_nv` fehlt in 89859 Einträgen, z.B. MastrNummer=ABR111222333123
- `Registergericht` fehlt in 95523 Einträgen, z.B. MastrNummer=ABR111222333123
- `Registergericht_nv` fehlt in 100000 Einträgen, z.B. MastrNummer=ABR111222333123
- `Registernummer` fehlt in 19211 Einträgen, z.B. MastrNummer=ABR900000008627
- `Registernummer_nv` fehlt in 100000 Einträgen, z.B. MastrNummer=ABR111222333123
- `Taetigkeitsende` fehlt in 99992 Einträgen, z.B. MastrNummer=ABR111222333123
- `Taetigkeitsende_nv` fehlt in 100000 Einträgen, z.B. MastrNummer=ABR111222333123
- `Telefon` fehlt in 92949 Einträgen, z.B. MastrNummer=ABR111222333123
- `Umsatzsteueridentifikationsnummer` fehlt in 95847 Einträgen, z.B. MastrNummer=ABR111222333123
- `Umsatzsteueridentifikationsnummer_nv` fehlt in 90662 Einträgen, z.B. MastrNummer=ABR111222333123
- `Webseite` fehlt in 97509 Einträgen, z.B. MastrNummer=ABR111222333123
- `Webseite_nv` fehlt in 90662 Einträgen, z.B. MastrNummer=ABR111222333123
```
```
Marktrollen.xml
- `BundesnetzagenturBetriebsnummer_nv` fehlt in 3447 Einträgen, z.B. MastrNummer=BVI900545854635EV
- `DatumLetzteAktualisierung` fehlt in 7371 Einträgen, z.B. MastrNummer=BVI900545854635EV
- `KontaktdatenMarktrolle` fehlt in 7479 Einträgen, z.B. MastrNummer=BVI900545854635EV
```
```
Netze.xml
- `Bezeichnung` fehlt in 803 Einträgen, z.B. MastrNummer=GNE904875516911
- `DatumLetzteAktualisierung` fehlt in 57 Einträgen, z.B. MastrNummer=GNE913909037773
- `GeschlossenesVerteilnetz` fehlt in 10 Einträgen, z.B. MastrNummer=GNE904875516911
- `KundenAngeschlossen` fehlt in 69 Einträgen, z.B. MastrNummer=GNE904875516911
```
```
AnlagenEegWind.xml
- `AnlagenkennzifferAnlagenregister` fehlt in 24175 Einträgen, z.B. EegMaStRNummer=EEG900001494788
- `AusschreibungZuschlag` fehlt in 901 Einträgen, z.B. EegMaStRNummer=EEG900199883449
- `PilotAnlage` fehlt in 26710 Einträgen, z.B. EegMaStRNummer=EEG900001494788
- `PrototypAnlage` fehlt in 5481 Einträgen, z.B. EegMaStRNummer=EEG900102733250
- `VerhaeltnisErtragsschaetzungReferenzertrag` fehlt in 29401 Einträgen, z.B. EegMaStRNummer=EEG900001494788
- `VerhaeltnisReferenzertragErtrag10Jahre` fehlt in 31137 Einträgen, z.B. EegMaStRNummer=EEG900001494788
- `VerhaeltnisReferenzertragErtrag15Jahre` fehlt in 31137 Einträgen, z.B. EegMaStRNummer=EEG900001494788
- `VerhaeltnisReferenzertragErtrag5Jahre` fehlt in 19113 Einträgen, z.B. EegMaStRNummer=EEG900007922711
```
```
AnlagenEegSolar_1.xml
- `AnlagenkennzifferAnlagenregister` fehlt in 99955 Einträgen, z.B. EegMaStRNummer=EEG900000009389
- `RegistrierungsnummerPvMeldeportal` fehlt in 75493 Einträgen, z.B. EegMaStRNummer=EEG900000016417
```
```
AnlagenEegBiomasse.xml
- `AnlagenkennzifferAnlagenregister` fehlt in 13103 Einträgen, z.B. EegMaStRNummer=EEG900007238652
- `AusschreibungZuschlag` fehlt in 1070 Einträgen, z.B. EegMaStRNummer=EEG900034100342
- `BiogasGaserzeugungskapazitaet` fehlt in 8876 Einträgen, z.B. EegMaStRNummer=EEG900012882178
- `BiomethanErstmaligerEinsatz` fehlt in 12568 Einträgen, z.B. EegMaStRNummer=EEG900007238652
- `Registrierungsdatum` fehlt in 1 Einträgen, z.B. EegMaStRNummer=EEG932023218028
```
```
AnlagenEegWasser.xml
- `AnlagenkennzifferAnlagenregister` fehlt in 5624 Einträgen, z.B. EegMaStRNummer=EEG900001696977
- `Registrierungsdatum` fehlt in 1 Einträgen, z.B. EegMaStRNummer=EEG943532633103
```
```
AnlagenEegGeoSolarthermieGrubenKlaerschlammDruckentspannung.xml
- `AnlagenkennzifferAnlagenregister` fehlt in 181 Einträgen, z.B. EegMaStRNummer=EEG901102403709
- `EegMastrNummer` fehlt in 184 Einträgen, z.B. EegMaStRNummer=EEG901102403709
AnlagenKwk.xml
```
```
- `AnlageBetriebsstatus` fehlt in 86985 Einträgen, z.B. KwkMastrNummer=KWK900003009119
- `Inbetriebnahmedatum` fehlt in 8773 Einträgen, z.B. KwkMastrNummer=KWK901420108287
- `Registrierungsdatum` fehlt in 8780 Einträgen, z.B. KwkMastrNummer=KWK906730599329
```
```
EinheitenGenehmigung.xml
- `Frist` fehlt in 18162 Einträgen, z.B. GenMastrNummer=SGE900003214940
```
```
EinheitenWind.xml
- `AnlagenbetreiberMastrNummer` fehlt in 397 Einträgen, z.B. EinheitMastrNummer=SEE900238810627
- `AuflageAbschaltungLeistungsbegrenzung` fehlt in 1756 Einträgen, z.B. EinheitMastrNummer=SEE900115838133
- `EegMaStRNummer` fehlt in 1830 Einträgen, z.B. EinheitMastrNummer=SEE900037558577
- `Hausnummer` fehlt in 30082 Einträgen, z.B. EinheitMastrNummer=SEE900002935310
- `Hausnummer_nv` fehlt in 778 Einträgen, z.B. EinheitMastrNummer=SEE900014803461
- `Kraftwerksnummer` fehlt in 32986 Einträgen, z.B. EinheitMastrNummer=SEE900002935310
- `Kuestenentfernung` fehlt in 31450 Einträgen, z.B. EinheitMastrNummer=SEE900002935310
- `LokationMaStRNummer` fehlt in 2231 Einträgen, z.B. EinheitMastrNummer=SEE900037558577
- `Nabenhoehe` fehlt in 737 Einträgen, z.B. EinheitMastrNummer=SEE900014803461
- `NameWindpark` fehlt in 32987 Einträgen, z.B. EinheitMastrNummer=SEE900002935310
- `Ort` fehlt in 1538 Einträgen, z.B. EinheitMastrNummer=SEE900019984141
- `Postleitzahl` fehlt in 1538 Einträgen, z.B. EinheitMastrNummer=SEE900019984141
- `Rotordurchmesser` fehlt in 392 Einträgen, z.B. EinheitMastrNummer=SEE900129535318
- `Typenbezeichnung` fehlt in 337 Einträgen, z.B. EinheitMastrNummer=SEE900129535318
- `Wassertiefe` fehlt in 31450 Einträgen, z.B. EinheitMastrNummer=SEE900002935310
- `Weic` fehlt in 31803 Einträgen, z.B. EinheitMastrNummer=SEE900002935310
```
```
EinheitenSolar_1.xml
- `AnlagenbetreiberMastrNummer` fehlt in 59 Einträgen, z.B. EinheitMastrNummer=SEE900088188204
- `EegMaStRNummer` fehlt in 377 Einträgen, z.B. EinheitMastrNummer=SEE900006249024
- `Hausnummer` fehlt in 92752 Einträgen, z.B. EinheitMastrNummer=SEE900000156564
- `Hausnummer_nv` fehlt in 91971 Einträgen, z.B. EinheitMastrNummer=SEE900000156564
- `Kraftwerksnummer` fehlt in 100000 Einträgen, z.B. EinheitMastrNummer=SEE900000156564
- `Weic` fehlt in 100000 Einträgen, z.B. EinheitMastrNummer=SEE900000156564
```
```
EinheitenBiomasse.xml
- `Hausnummer` fehlt in 4946 Einträgen, z.B. EinheitMastrNummer=SEE900000911668
- `Hausnummer_nv` fehlt in 831 Einträgen, z.B. EinheitMastrNummer=SEE900000911668
- `Kraftwerksnummer` fehlt in 21089 Einträgen, z.B. EinheitMastrNummer=SEE900000911668
- `LokationMaStRNummer` fehlt in 518 Einträgen, z.B. EinheitMastrNummer=SEE900004080537
- `Nettonennleistung` fehlt in 6 Einträgen, z.B. EinheitMastrNummer=SEE907541917340
- `Ort` fehlt in 2 Einträgen, z.B. EinheitMastrNummer=SEE930479566687
- `Postleitzahl` fehlt in 2 Einträgen, z.B. EinheitMastrNummer=SEE930479566687
- `Technologie` fehlt in 3 Einträgen, z.B. EinheitMastrNummer=SEE907541917340
- `Weic` fehlt in 21064 Einträgen, z.B. EinheitMastrNummer=SEE900000911668
```
```
EinheitenWasser.xml
- `AnlagenbetreiberMastrNummer` fehlt in 22 Einträgen, z.B. EinheitMastrNummer=SEE904155331845
- `EegMaStRNummer` fehlt in 1875 Einträgen, z.B. EinheitMastrNummer=SEE900012607467
- `Hausnummer` fehlt in 4580 Einträgen, z.B. EinheitMastrNummer=SEE900012607467
- `Hausnummer_nv` fehlt in 3818 Einträgen, z.B. EinheitMastrNummer=SEE900012607467
- `Kraftwerksnummer` fehlt in 8421 Einträgen, z.B. EinheitMastrNummer=SEE900012607467
- `LokationMaStRNummer` fehlt in 58 Einträgen, z.B. EinheitMastrNummer=SEE901721465997
- `Nettonennleistung` fehlt in 1 Einträgen, z.B. EinheitMastrNummer=SEE930688632715
- `Weic` fehlt in 8380 Einträgen, z.B. EinheitMastrNummer=SEE900012607467
```
```
EinheitenGeoSolarthermieGrubenKlaerschlammDruckentspannung.xml
- `Hausnummer` fehlt in 184 Einträgen, z.B. EinheitMastrNummer=SEE900331826662
- `Hausnummer_nv` fehlt in 74 Einträgen, z.B. EinheitMastrNummer=SEE901437097857
- `Kraftwerksnummer` fehlt in 455 Einträgen, z.B. EinheitMastrNummer=SEE900331826662
- `LokationMaStRNummer` fehlt in 15 Einträgen, z.B. EinheitMastrNummer=SEE900842254756
- `Weic` fehlt in 455 Einträgen, z.B. EinheitMastrNummer=SEE900331826662
```
```
EinheitenVerbrennung.xml
- `AnlagenbetreiberMastrNummer` fehlt in 1193 Einträgen, z.B. EinheitMastrNummer=SEE900073920406
- `AnzeigeEinerStilllegung` fehlt in 76355 Einträgen, z.B. EinheitMastrNummer=SEE900001203007
- `ArtDerStilllegung` fehlt in 77102 Einträgen, z.B. EinheitMastrNummer=SEE900001203007
- `Hauptbrennstoff` fehlt in 2479 Einträgen, z.B. EinheitMastrNummer=SEE900005568744
- `Hausnummer` fehlt in 57880 Einträgen, z.B. EinheitMastrNummer=SEE900001203007
- `Hausnummer_nv` fehlt in 56903 Einträgen, z.B. EinheitMastrNummer=SEE900001203007
- `Kraftwerksnummer` fehlt in 76565 Einträgen, z.B. EinheitMastrNummer=SEE900001203007
- `LokationMaStRNummer` fehlt in 1748 Einträgen, z.B. EinheitMastrNummer=SEE900073920406
- `NameKraftwerk` fehlt in 45142 Einträgen, z.B. EinheitMastrNummer=SEE900002061170
- `Nettonennleistung` fehlt in 28 Einträgen, z.B. EinheitMastrNummer=SEE903722032437
- `Weic` fehlt in 76687 Einträgen, z.B. EinheitMastrNummer=SEE900001203007
- `WeitereBrennstoffe` fehlt in 74749 Einträgen, z.B. EinheitMastrNummer=SEE900001203007
- `WeitererHauptbrennstoff` fehlt in 68892 Einträgen, z.B. EinheitMastrNummer=SEE900001203007
```
```
EinheitenKernkraft.xml
- `AnlagenbetreiberMastrNummer` fehlt in 1 Einträgen, z.B. EinheitMastrNummer=SEE927528071629
- `Hausnummer` fehlt in 4 Einträgen, z.B. EinheitMastrNummer=SEE930752846949
- `Kraftwerksnummer` fehlt in 4 Einträgen, z.B. EinheitMastrNummer=SEE930752846949
- `LokationMaStRNummer` fehlt in 1 Einträgen, z.B. EinheitMastrNummer=SEE927528071629
- `Weic` fehlt in 3 Einträgen, z.B. EinheitMastrNummer=SEE930752846949
```
```
EinheitenStromSpeicher_1.xml
- `AnlagenbetreiberMastrNummer` fehlt in 135 Einträgen, z.B. EinheitMastrNummer=SEE900357122236
- `EegMaStRNummer` fehlt in 3397 Einträgen, z.B. EinheitMastrNummer=SEE900004095760
- `Hausnummer` fehlt in 99719 Einträgen, z.B. EinheitMastrNummer=SEE900000066023
- `Hausnummer_nv` fehlt in 99663 Einträgen, z.B. EinheitMastrNummer=SEE900000066023
- `Kraftwerksnummer` fehlt in 99966 Einträgen, z.B. EinheitMastrNummer=SEE900000066023
- `LokationMaStRNummer` fehlt in 1367 Einträgen, z.B. EinheitMastrNummer=SEE900075693025
- `Nettonennleistung` fehlt in 3 Einträgen, z.B. EinheitMastrNummer=SEE910814993941
- `Weic` fehlt in 99977 Einträgen, z.B. EinheitMastrNummer=SEE900000066023
```
```
EinheitenStromVerbraucher.xml
- `Hausnummer` fehlt in 69 Einträgen, z.B. EinheitMastrNummer=SVE900066805970
- `LokationMaStRNummer` fehlt in 47 Einträgen, z.B. EinheitMastrNummer=SVE900530977578
```
```
EinheitenGasErzeuger.xml
- `Erzeugungsleistung` fehlt in 53 Einträgen, z.B. EinheitMastrNummer=GEE900074558080
- `Hausnummer` fehlt in 108 Einträgen, z.B. EinheitMastrNummer=GEE900074558080
- `Hausnummer_nv` fehlt in 53 Einträgen, z.B. EinheitMastrNummer=GEE900074558080
- `LokationMaStRNummer` fehlt in 7 Einträgen, z.B. EinheitMastrNummer=GEE903148821912
- `SpeicherMaStRNummer` fehlt in 254 Einträgen, z.B. EinheitMastrNummer=GEE900126760591
```
```
EinheitenGasSpeicher.xml
- `Hausnummer` fehlt in 6 Einträgen, z.B. EinheitMastrNummer=GEE941903701734
- `Weic` fehlt in 30 Einträgen, z.B. EinheitMastrNummer=GEE905950524425
- `Weic_nv` fehlt in 53 Einträgen, z.B. EinheitMastrNummer=GEE900074558080
```
```
EinheitenGasverbraucher.xml
- `AnlagenbetreiberMastrNummer` fehlt in 9 Einträgen, z.B. EinheitMastrNummer=GVE916010877795
- `Hausnummer` fehlt in 55 Einträgen, z.B. EinheitMastrNummer=GVE904101321760
- `LokationMaStRNummer` fehlt in 27 Einträgen, z.B. EinheitMastrNummer=GVE901033793832
- `NameGasverbrauchseinheit` fehlt in 718 Einträgen, z.B. EinheitMastrNummer=GVE900211727539
```

## Unterschiede in der Schreibweise

Eineiten


...

# Probleme mit dem Datenformat

UTF-16
XML

<aside>Interessant hierbei ist, dass der Datenexport über eine Stunde in Anspruch zu nehmen scheint -- zu erkennen an der Differenz der Zeitstempel zwischen dem ersten Eintrag um 01:00 und dem letzten um 02:09.

# TODO

- Lächerlicher Downloadlink




[broken-ex1]: https://archive.is/EYyIQ
[broken-ex2]: https://archive.is/kyxmX
[export-archived]: https://s3.eu-central-1.wasabisys.com/mastr-backup/Gesamtdatenexport_20220203__840cfde7b693453982d28db827025ff0.zip
[export-doc-archived]: https://web.archive.org/web/20211222091956/https://www.marktstammdatenregister.de/MaStRHilfe/files/gesamtdatenexport/Dokumentation%20MaStR%20Gesamtdatenexport.pdf
[export-doc-live]: https://www.marktstammdatenregister.de/MaStRHilfe/files/gesamtdatenexport/Dokumentation%20MaStR%20Gesamtdatenexport.pdf
[export-download]: https://www.marktstammdatenregister.de/MaStR/Datendownload
[mastr]: https://www.marktstammdatenregister.de/MaStR/
