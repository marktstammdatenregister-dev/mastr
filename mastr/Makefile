# Usage:
#
#     make download
#     make Marktstammdatenregister.db.br

# Download the latest export.
Gesamtdatenexport.zip:
	curl -sSL https://www.marktstammdatenregister.de/MaStR/Datendownload | \
		pup 'a[href^="https://download.marktstammdatenregister.de/Gesamtdatenexport"][href$$=".zip"] attr{href}' | \
		xargs axel --output=$@

.PHONY: download
download: Gesamtdatenexport.zip

# Run 'make download' to populate these variables.
EinheitenSolar_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'EinheitenSolar_')
EinheitenSolar_csv := $(EinheitenSolar_xml:%.xml=%.csv)

EinheitenWind_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'EinheitenWind')
EinheitenWind_csv := $(EinheitenWind_xml:%.xml=%.csv)

EinheitenWasser_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'EinheitenWasser')
EinheitenWasser_csv := $(EinheitenWasser_xml:%.xml=%.csv)

EinheitenStromSpeicher_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'EinheitenStromSpeicher_')
EinheitenStromSpeicher_csv := $(EinheitenStromSpeicher_xml:%.xml=%.csv)

AnlagenEegSolar_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'AnlagenEegSolar_')
AnlagenEegSolar_csv := $(AnlagenEegSolar_xml:%.xml=%.csv)

AnlagenEegWind_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'AnlagenEegWind')
AnlagenEegWind_csv := $(AnlagenEegWind_xml:%.xml=%.csv)

AnlagenEegWasser_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'AnlagenEegWasser')
AnlagenEegWasser_csv := $(AnlagenEegWasser_xml:%.xml=%.csv)

AnlagenStromSpeicher_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'AnlagenStromSpeicher_')
AnlagenStromSpeicher_csv := $(AnlagenStromSpeicher_xml:%.xml=%.csv)

Lokationen_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'Lokationen_')
Lokationen_csv := $(Lokationen_xml:%.xml=%.csv)

Marktakteure_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'Marktakteure_')
Marktakteure_csv := $(Marktakteure_xml:%.xml=%.csv)

Katalogkategorien_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'Katalogkategorien')
Katalogkategorien_csv := $(Katalogkategorien_xml:%.xml=%.csv)

Katalogwerte_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'Katalogwerte')
Katalogwerte_csv := $(Katalogwerte_xml:%.xml=%.csv)

Marktstammdatenregister_sql := EinheitenSolar-csv.gen.sql EinheitenWind-csv.gen.sql EinheitenWasser-csv.gen.sql EinheitenStromSpeicher-csv.gen.sql AnlagenEegSolar-csv.gen.sql AnlagenEegWind-csv.gen.sql AnlagenEegWasser-csv.gen.sql AnlagenStromSpeicher-csv.gen.sql Lokationen-csv.gen.sql Marktakteure-csv.gen.sql Katalogkategorien-csv.gen.sql Katalogwerte-csv.gen.sql

%.gen.xsd: %.yml XSD.mustache
	mustache $< XSD.mustache >$@

%.gen.sql: %.yml SQLite.mustache
	mustache $< SQLite.mustache >$@

%.gen.json: %.yml
	yj <$< >$@

xml2csv: xml2csv.go
	go build -o $@ $<

EinheitenSolar_%.csv: Gesamtdatenexport.zip EinheitenSolar.gen.xsd EinheitenSolar.gen.json xml2csv
	unzip -qo $< $(subst ,,$(@:%.csv=%.xml))
	xmllint --schema EinheitenSolar.gen.xsd --noout $(@:%.csv=%.xml)
	iconv -f utf-16 -t utf-8 <$(@:%.csv=%.xml) | ./xml2csv -descriptor EinheitenSolar.gen.json >$@

EinheitenWind.csv: Gesamtdatenexport.zip EinheitenWind.gen.xsd EinheitenWind.gen.json xml2csv
	unzip -qo $< $(subst ,,$(@:%.csv=%.xml))
	xmllint --schema EinheitenWind.gen.xsd --noout $(@:%.csv=%.xml)
	iconv -f utf-16 -t utf-8 <$(@:%.csv=%.xml) | ./xml2csv -descriptor EinheitenWind.gen.json >$@

EinheitenWasser.csv: Gesamtdatenexport.zip EinheitenWasser.gen.xsd EinheitenWasser.gen.json xml2csv
	unzip -qo $< $(subst ,,$(@:%.csv=%.xml))
	xmllint --schema EinheitenWasser.gen.xsd --noout $(@:%.csv=%.xml)
	iconv -f utf-16 -t utf-8 <$(@:%.csv=%.xml) | ./xml2csv -descriptor EinheitenWasser.gen.json >$@

EinheitenStromSpeicher_%.csv: Gesamtdatenexport.zip EinheitenStromSpeicher.gen.xsd EinheitenStromSpeicher.gen.json xml2csv
	unzip -qo $< $(subst ,,$(@:%.csv=%.xml))
	xmllint --schema EinheitenStromSpeicher.gen.xsd --noout $(@:%.csv=%.xml)
	iconv -f utf-16 -t utf-8 <$(@:%.csv=%.xml) | ./xml2csv -descriptor EinheitenStromSpeicher.gen.json >$@

AnlagenEegSolar_%.csv: Gesamtdatenexport.zip AnlagenEegSolar.gen.xsd AnlagenEegSolar.gen.json xml2csv
	unzip -qo $< $(subst ,,$(@:%.csv=%.xml))
	xmllint --schema AnlagenEegSolar.gen.xsd --noout $(@:%.csv=%.xml)
	iconv -f utf-16 -t utf-8 <$(@:%.csv=%.xml) | ./xml2csv -descriptor AnlagenEegSolar.gen.json >$@

AnlagenEegWind.csv: Gesamtdatenexport.zip AnlagenEegWind.gen.xsd AnlagenEegWind.gen.json xml2csv
	unzip -qo $< $(subst ,,$(@:%.csv=%.xml))
	xmllint --schema AnlagenEegWind.gen.xsd --noout $(@:%.csv=%.xml)
	iconv -f utf-16 -t utf-8 <$(@:%.csv=%.xml) | ./xml2csv -descriptor AnlagenEegWind.gen.json >$@

AnlagenEegWasser.csv: Gesamtdatenexport.zip AnlagenEegWasser.gen.xsd AnlagenEegWasser.gen.json xml2csv
	unzip -qo $< $(subst ,,$(@:%.csv=%.xml))
	xmllint --schema AnlagenEegWasser.gen.xsd --noout $(@:%.csv=%.xml)
	iconv -f utf-16 -t utf-8 <$(@:%.csv=%.xml) | ./xml2csv -descriptor AnlagenEegWasser.gen.json >$@

AnlagenStromSpeicher_%.csv: Gesamtdatenexport.zip AnlagenStromSpeicher.gen.xsd AnlagenStromSpeicher.gen.json xml2csv
	unzip -qo $< $(subst ,,$(@:%.csv=%.xml))
	xmllint --schema AnlagenStromSpeicher.gen.xsd --noout $(@:%.csv=%.xml)
	iconv -f utf-16 -t utf-8 <$(@:%.csv=%.xml) | ./xml2csv -descriptor AnlagenStromSpeicher.gen.json >$@

Lokationen_%.csv: Gesamtdatenexport.zip Lokationen.gen.xsd Lokationen.gen.json xml2csv
	unzip -qo $< $(subst ,,$(@:%.csv=%.xml))
	xmllint --schema Lokationen.gen.xsd --noout $(@:%.csv=%.xml)
	iconv -f utf-16 -t utf-8 <$(@:%.csv=%.xml) | ./xml2csv -descriptor Lokationen.gen.json >$@

Marktakteure_%.csv: Gesamtdatenexport.zip Marktakteure.gen.xsd Marktakteure.gen.json xml2csv
	unzip -qo $< $(subst ,,$(@:%.csv=%.xml))
	xmllint --schema Marktakteure.gen.xsd --noout $(@:%.csv=%.xml)
	iconv -f utf-16 -t utf-8 <$(@:%.csv=%.xml) | ./xml2csv -descriptor Marktakteure.gen.json >$@

Katalogkategorien.csv: Gesamtdatenexport.zip Katalogkategorien.gen.xsd Katalogkategorien.gen.json xml2csv
	unzip -qo $< $(subst ,,$(@:%.csv=%.xml))
	xmllint --schema Katalogkategorien.gen.xsd --noout $(@:%.csv=%.xml)
	iconv -f utf-16 -t utf-8 <$(@:%.csv=%.xml) | ./xml2csv -descriptor Katalogkategorien.gen.json >$@

Katalogwerte.csv: Gesamtdatenexport.zip Katalogwerte.gen.xsd Katalogwerte.gen.json xml2csv
	unzip -qo $< $(subst ,,$(@:%.csv=%.xml))
	xmllint --schema Katalogwerte.gen.xsd --noout $(@:%.csv=%.xml)
	iconv -f utf-16 -t utf-8 <$(@:%.csv=%.xml) | ./xml2csv -descriptor Katalogwerte.gen.json >$@

EinheitenSolar-csv.gen.sql: EinheitenSolar.gen.sql $(EinheitenSolar_csv)
	cp EinheitenSolar.gen.sql $@
	echo '.mode csv' >>$@
	printf "%s\0" $(EinheitenSolar_csv) | \
		xargs -0 --no-run-if-empty -I'{}' sh -c 'echo ".import {} EinheitenSolar --skip 1" >>$@'

EinheitenWind-csv.gen.sql: EinheitenWind.gen.sql $(EinheitenWind_csv)
	cp EinheitenWind.gen.sql $@
	echo '.mode csv' >>$@
	printf "%s\0" $(EinheitenWind_csv) | \
		xargs -0 --no-run-if-empty -I'{}' sh -c 'echo ".import {} EinheitenWind --skip 1" >>$@'

EinheitenWasser-csv.gen.sql: EinheitenWasser.gen.sql $(EinheitenWasser_csv)
	cp EinheitenWasser.gen.sql $@
	echo '.mode csv' >>$@
	printf "%s\0" $(EinheitenWasser_csv) | \
		xargs -0 --no-run-if-empty -I'{}' sh -c 'echo ".import {} EinheitenWasser --skip 1" >>$@'

EinheitenStromSpeicher-csv.gen.sql: EinheitenStromSpeicher.gen.sql $(EinheitenStromSpeicher_csv)
	cp EinheitenStromSpeicher.gen.sql $@
	echo '.mode csv' >>$@
	printf "%s\0" $(EinheitenStromSpeicher_csv) | \
		xargs -0 --no-run-if-empty -I'{}' sh -c 'echo ".import {} EinheitenStromSpeicher --skip 1" >>$@'

AnlagenEegSolar-csv.gen.sql: AnlagenEegSolar.gen.sql $(AnlagenEegSolar_csv)
	cp AnlagenEegSolar.gen.sql $@
	echo '.mode csv' >>$@
	printf "%s\0" $(AnlagenEegSolar_csv) | \
		xargs -0 --no-run-if-empty -I'{}' sh -c 'echo ".import {} AnlagenEegSolar --skip 1" >>$@'

AnlagenEegWind-csv.gen.sql: AnlagenEegWind.gen.sql $(AnlagenEegWind_csv)
	cp AnlagenEegWind.gen.sql $@
	echo '.mode csv' >>$@
	printf "%s\0" $(AnlagenEegWind_csv) | \
		xargs -0 --no-run-if-empty -I'{}' sh -c 'echo ".import {} AnlagenEegWind --skip 1" >>$@'

AnlagenEegWasser-csv.gen.sql: AnlagenEegWasser.gen.sql $(AnlagenEegWasser_csv)
	cp AnlagenEegWasser.gen.sql $@
	echo '.mode csv' >>$@
	printf "%s\0" $(AnlagenEegWasser_csv) | \
		xargs -0 --no-run-if-empty -I'{}' sh -c 'echo ".import {} AnlagenEegWasser --skip 1" >>$@'

AnlagenStromSpeicher-csv.gen.sql: AnlagenStromSpeicher.gen.sql $(AnlagenStromSpeicher_csv)
	cp AnlagenStromSpeicher.gen.sql $@
	echo '.mode csv' >>$@
	printf "%s\0" $(AnlagenStromSpeicher_csv) | \
		xargs -0 --no-run-if-empty -I'{}' sh -c 'echo ".import {} AnlagenStromSpeicher --skip 1" >>$@'

Lokationen-csv.gen.sql: Lokationen.gen.sql $(Lokationen_csv)
	cp Lokationen.gen.sql $@
	echo '.mode csv' >>$@
	printf "%s\0" $(Lokationen_csv) | \
		xargs -0 --no-run-if-empty -I'{}' sh -c 'echo ".import {} Lokationen --skip 1" >>$@'

Marktakteure-csv.gen.sql: Marktakteure.gen.sql $(Marktakteure_csv)
	cp Marktakteure.gen.sql $@
	echo '.mode csv' >>$@
	printf "%s\0" $(Marktakteure_csv) | \
		xargs -0 --no-run-if-empty -I'{}' sh -c 'echo ".import {} Marktakteure --skip 1" >>$@'

Katalogkategorien-csv.gen.sql: Katalogkategorien.gen.sql $(Katalogkategorien_csv)
	cp Katalogkategorien.gen.sql $@
	echo '.mode csv' >>$@
	printf "%s\0" $(Katalogkategorien_csv) | \
		xargs -0 --no-run-if-empty -I'{}' sh -c 'echo ".import {} Katalogkategorien --skip 1" >>$@'

Katalogwerte-csv.gen.sql: Katalogwerte.gen.sql $(Katalogwerte_csv)
	cp Katalogwerte.gen.sql $@
	echo '.mode csv' >>$@
	printf "%s\0" $(Katalogwerte_csv) | \
		xargs -0 --no-run-if-empty -I'{}' sh -c 'echo ".import {} Katalogwerte --skip 1" >>$@'

Marktstammdatenregister.db: $(Marktstammdatenregister_sql)
	rm -f $@
	cat $(Marktstammdatenregister_sql) | sqlite3 $@
	sqlite3 $@ <<<'CREATE INDEX idx_EinheitenSolar_Ort_Nettonennleistung on EinheitenSolar(Ort, Nettonennleistung);'
	sqlite3 $@ <<<'CREATE INDEX idx_EinheitenSolar_Inbetriebnahmedatum_Nettonennleistung on EinheitenSolar(Inbetriebnahmedatum, Nettonennleistung);'
	sqlite3 $@ <<<'CREATE INDEX idx_EinheitenSolar_DatumEndgueltigeStilllegung_Nettonennleistung on EinheitenSolar(DatumEndgueltigeStilllegung, Nettonennleistung);'
	sqlite3 $@ <<<'CREATE INDEX idx_EinheitenSolar_Bundesland on EinheitenSolar(Bundesland);'
	sqlite3 $@ <<<'CREATE INDEX idx_EinheitenSolar_Landkreis on EinheitenSolar(Landkreis);'
	sqlite3 $@ <<<'CREATE INDEX idx_EinheitenSolar_Nutzungsbereich on EinheitenSolar(Nutzungsbereich);'
	sqlite3 $@ <<<'CREATE INDEX idx_EinheitenSolar_AnlagenbetreiberMastrNummer on EinheitenSolar(AnlagenbetreiberMastrNummer);'
	sqlite3 $@ <<<'VACUUM; ANALYZE;'

Marktstammdatenregister.db.br: Marktstammdatenregister.db
	brotli -4 --keep --force --output=$@ $<
	touch $@

.PHONY: clean
clean:
	rm -f *.gen.*
	rm -f $(EinheitenSolar_xml) $(EinheitenSolar_csv)
	rm -f $(EinheitenWind_xml) $(EinheitenWind_csv)
	rm -f $(EinheitenWasser_xml) $(EinheitenWasser_csv)
	rm -f $(EinheitenStromSpeicher_xml) $(EinheitenStromSpeicher_csv)
	rm -f $(AnlagenEegSolar_xml) $(AnlagenEegSolar_csv)
	rm -f $(AnlagenEegWind_xml) $(AnlagenEegWind_csv)
	rm -f $(AnlagenEegWasser_xml) $(AnlagenEegWasser_csv)
	rm -f $(AnlagenStromSpeicher_xml) $(AnlagenStromSpeicher_csv)
	rm -f $(Lokationen_xml) $(Lokationen_csv)
	rm -f $(Katalogwerte_xml) $(Katalogwerte_csv)
	rm -f $(Katalogkategorien_xml) $(Katalogkategorien_csv)
	rm -f $(Marktakteure_xml) $(Marktakteure_csv)
	rm -f $(Marktstammdatenregister_sql)
	rm -f Marktstammdatenregister.db{,.br}
