# Usage:
#
#     make download
#     make Marktstammdatenregister.db.br

SHELL := bash

# Download the latest export.
Gesamtdatenexport.zip:
	curl -sSL https://www.marktstammdatenregister.de/MaStR/Datendownload | \
		pup 'a[href^="https://download.marktstammdatenregister.de/Gesamtdatenexport"][href$$=".zip"] attr{href}' | \
		xargs axel --quiet --output=$@

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

EinheitenBiomasse_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'EinheitenBiomasse')
EinheitenBiomasse_csv := $(EinheitenBiomasse_xml:%.xml=%.csv)

AnlagenEegSolar_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'AnlagenEegSolar_')
AnlagenEegSolar_csv := $(AnlagenEegSolar_xml:%.xml=%.csv)

AnlagenEegWind_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'AnlagenEegWind')
AnlagenEegWind_csv := $(AnlagenEegWind_xml:%.xml=%.csv)

AnlagenEegWasser_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'AnlagenEegWasser')
AnlagenEegWasser_csv := $(AnlagenEegWasser_xml:%.xml=%.csv)

AnlagenStromSpeicher_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'AnlagenStromSpeicher_')
AnlagenStromSpeicher_csv := $(AnlagenStromSpeicher_xml:%.xml=%.csv)

AnlagenEegBiomasse_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'AnlagenEegBiomasse')
AnlagenEegBiomasse_csv := $(AnlagenEegBiomasse_xml:%.xml=%.csv)

Lokationen_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'Lokationen_')
Lokationen_csv := $(Lokationen_xml:%.xml=%.csv)

Marktakteure_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'Marktakteure_')
Marktakteure_csv := $(Marktakteure_xml:%.xml=%.csv)

Katalogkategorien_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'Katalogkategorien')
Katalogkategorien_csv := $(Katalogkategorien_xml:%.xml=%.csv)

Katalogwerte_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'Katalogwerte')
Katalogwerte_csv := $(Katalogwerte_xml:%.xml=%.csv)

Marktstammdatenregister_sql := \
Lokationen-csv.gen.sql \
Katalogkategorien-csv.gen.sql \
Katalogwerte-csv.gen.sql \
Marktakteure-csv.gen.sql \
AnlagenEegBiomasse-csv.gen.sql \
AnlagenEegSolar-csv.gen.sql \
AnlagenEegWasser-csv.gen.sql \
AnlagenEegWind-csv.gen.sql \
AnlagenStromSpeicher-csv.gen.sql \
EinheitenBiomasse-csv.gen.sql \
EinheitenSolar-csv.gen.sql \
EinheitenStromSpeicher-csv.gen.sql \
EinheitenWasser-csv.gen.sql \
EinheitenWind-csv.gen.sql

%.gen.xsd: %.yaml XSD.mustache
	mustache $< XSD.mustache >$@

%.gen.sql: %.yaml PostgreSQL.mustache
	mustache $< PostgreSQL.mustache >$@

xml2csv: xml2csv.go
	go build -o $@ $<

EinheitenSolar_%.csv: Gesamtdatenexport.zip EinheitenSolar.gen.xsd EinheitenSolar.yaml xml2csv
	unzip -qo $< $(subst ,,$(@:%.csv=%.xml))
	xmllint --schema EinheitenSolar.gen.xsd --noout $(@:%.csv=%.xml)
	iconv -f utf-16 -t utf-8 <$(@:%.csv=%.xml) | ./xml2csv -descriptor EinheitenSolar.yaml >$@

EinheitenWind.csv: Gesamtdatenexport.zip EinheitenWind.gen.xsd EinheitenWind.yaml xml2csv
	unzip -qo $< $(subst ,,$(@:%.csv=%.xml))
	xmllint --schema EinheitenWind.gen.xsd --noout $(@:%.csv=%.xml)
	iconv -f utf-16 -t utf-8 <$(@:%.csv=%.xml) | ./xml2csv -descriptor EinheitenWind.yaml >$@

EinheitenWasser.csv: Gesamtdatenexport.zip EinheitenWasser.gen.xsd EinheitenWasser.yaml xml2csv
	unzip -qo $< $(subst ,,$(@:%.csv=%.xml))
	xmllint --schema EinheitenWasser.gen.xsd --noout $(@:%.csv=%.xml)
	iconv -f utf-16 -t utf-8 <$(@:%.csv=%.xml) | ./xml2csv -descriptor EinheitenWasser.yaml >$@

EinheitenStromSpeicher_%.csv: Gesamtdatenexport.zip EinheitenStromSpeicher.gen.xsd EinheitenStromSpeicher.yaml xml2csv
	unzip -qo $< $(subst ,,$(@:%.csv=%.xml))
	xmllint --schema EinheitenStromSpeicher.gen.xsd --noout $(@:%.csv=%.xml)
	iconv -f utf-16 -t utf-8 <$(@:%.csv=%.xml) | ./xml2csv -descriptor EinheitenStromSpeicher.yaml >$@

EinheitenBiomasse.csv: Gesamtdatenexport.zip EinheitenBiomasse.gen.xsd EinheitenBiomasse.yaml xml2csv
	unzip -qo $< $(subst ,,$(@:%.csv=%.xml))
	xmllint --schema EinheitenBiomasse.gen.xsd --noout $(@:%.csv=%.xml)
	iconv -f utf-16 -t utf-8 <$(@:%.csv=%.xml) | ./xml2csv -descriptor EinheitenBiomasse.yaml >$@

AnlagenEegSolar_%.csv: Gesamtdatenexport.zip AnlagenEegSolar.gen.xsd AnlagenEegSolar.yaml xml2csv
	unzip -qo $< $(subst ,,$(@:%.csv=%.xml))
	xmllint --schema AnlagenEegSolar.gen.xsd --noout $(@:%.csv=%.xml)
	iconv -f utf-16 -t utf-8 <$(@:%.csv=%.xml) | ./xml2csv -descriptor AnlagenEegSolar.yaml >$@

AnlagenEegWind.csv: Gesamtdatenexport.zip AnlagenEegWind.gen.xsd AnlagenEegWind.yaml xml2csv
	unzip -qo $< $(subst ,,$(@:%.csv=%.xml))
	xmllint --schema AnlagenEegWind.gen.xsd --noout $(@:%.csv=%.xml)
	iconv -f utf-16 -t utf-8 <$(@:%.csv=%.xml) | ./xml2csv -descriptor AnlagenEegWind.yaml >$@

AnlagenEegWasser.csv: Gesamtdatenexport.zip AnlagenEegWasser.gen.xsd AnlagenEegWasser.yaml xml2csv
	unzip -qo $< $(subst ,,$(@:%.csv=%.xml))
	xmllint --schema AnlagenEegWasser.gen.xsd --noout $(@:%.csv=%.xml)
	iconv -f utf-16 -t utf-8 <$(@:%.csv=%.xml) | ./xml2csv -descriptor AnlagenEegWasser.yaml >$@

AnlagenStromSpeicher_%.csv: Gesamtdatenexport.zip AnlagenStromSpeicher.gen.xsd AnlagenStromSpeicher.yaml xml2csv
	unzip -qo $< $(subst ,,$(@:%.csv=%.xml))
	xmllint --schema AnlagenStromSpeicher.gen.xsd --noout $(@:%.csv=%.xml)
	iconv -f utf-16 -t utf-8 <$(@:%.csv=%.xml) | ./xml2csv -descriptor AnlagenStromSpeicher.yaml >$@

AnlagenEegBiomasse.csv: Gesamtdatenexport.zip AnlagenEegBiomasse.gen.xsd AnlagenEegBiomasse.yaml xml2csv
	unzip -qo $< $(subst ,,$(@:%.csv=%.xml))
	xmllint --schema AnlagenEegBiomasse.gen.xsd --noout $(@:%.csv=%.xml)
	iconv -f utf-16 -t utf-8 <$(@:%.csv=%.xml) | ./xml2csv -descriptor AnlagenEegBiomasse.yaml >$@

Lokationen_%.csv: Gesamtdatenexport.zip Lokationen.gen.xsd Lokationen.yaml xml2csv
	unzip -qo $< $(subst ,,$(@:%.csv=%.xml))
	xmllint --schema Lokationen.gen.xsd --noout $(@:%.csv=%.xml)
	iconv -f utf-16 -t utf-8 <$(@:%.csv=%.xml) | ./xml2csv -descriptor Lokationen.yaml >$@

Marktakteure_%.csv: Gesamtdatenexport.zip Marktakteure.gen.xsd Marktakteure.yaml xml2csv
	unzip -qo $< $(subst ,,$(@:%.csv=%.xml))
	xmllint --schema Marktakteure.gen.xsd --noout $(@:%.csv=%.xml)
	iconv -f utf-16 -t utf-8 <$(@:%.csv=%.xml) | ./xml2csv -descriptor Marktakteure.yaml >$@

Katalogkategorien.csv: Gesamtdatenexport.zip Katalogkategorien.gen.xsd Katalogkategorien.yaml xml2csv
	unzip -qo $< $(subst ,,$(@:%.csv=%.xml))
	xmllint --schema Katalogkategorien.gen.xsd --noout $(@:%.csv=%.xml)
	iconv -f utf-16 -t utf-8 <$(@:%.csv=%.xml) | ./xml2csv -descriptor Katalogkategorien.yaml >$@

Katalogwerte.csv: Gesamtdatenexport.zip Katalogwerte.gen.xsd Katalogwerte.yaml xml2csv
	unzip -qo $< $(subst ,,$(@:%.csv=%.xml))
	xmllint --schema Katalogwerte.gen.xsd --noout $(@:%.csv=%.xml)
	iconv -f utf-16 -t utf-8 <$(@:%.csv=%.xml) | ./xml2csv -descriptor Katalogwerte.yaml >$@

EinheitenSolar-csv.gen.sql: EinheitenSolar.gen.sql $(EinheitenSolar_csv)
	echo 'BEGIN;' >$@
	cat EinheitenSolar.gen.sql >>$@
	for f in $(EinheitenSolar_csv); do \
		printf '\COPY "EinheitenSolar" FROM %s CSV;\n' "'$$f'" >>"$@"; \
	done
	echo 'COMMIT;' >>$@

EinheitenWind-csv.gen.sql: EinheitenWind.gen.sql $(EinheitenWind_csv)
	echo 'BEGIN;' >$@
	cat EinheitenWind.gen.sql >>$@
	for f in $(EinheitenWind_csv); do \
		printf '\COPY "EinheitenWind" FROM %s CSV;\n' "'$$f'" >>"$@"; \
	done
	echo 'COMMIT;' >>$@

EinheitenWasser-csv.gen.sql: EinheitenWasser.gen.sql $(EinheitenWasser_csv)
	echo 'BEGIN;' >$@
	cat EinheitenWasser.gen.sql >>$@
	for f in $(EinheitenWasser_csv); do \
		printf '\COPY "EinheitenWasser" FROM %s CSV;\n' "'$$f'" >>"$@"; \
	done
	echo 'COMMIT;' >>$@

EinheitenStromSpeicher-csv.gen.sql: EinheitenStromSpeicher.gen.sql $(EinheitenStromSpeicher_csv)
	echo 'BEGIN;' >$@
	cat EinheitenStromSpeicher.gen.sql >>$@
	for f in $(EinheitenStromSpeicher_csv); do \
		printf '\COPY "EinheitenStromSpeicher" FROM %s CSV;\n' "'$$f'" >>"$@"; \
	done
	echo 'COMMIT;' >>$@

EinheitenBiomasse-csv.gen.sql: EinheitenBiomasse.gen.sql $(EinheitenBiomasse_csv)
	echo 'BEGIN;' >$@
	cat EinheitenBiomasse.gen.sql >>$@
	for f in $(EinheitenBiomasse_csv); do \
		printf '\COPY "EinheitenBiomasse" FROM %s CSV;\n' "'$$f'" >>"$@"; \
	done
	echo 'COMMIT;' >>$@

AnlagenEegSolar-csv.gen.sql: AnlagenEegSolar.gen.sql $(AnlagenEegSolar_csv)
	echo 'BEGIN;' >$@
	cat AnlagenEegSolar.gen.sql >>$@
	for f in $(AnlagenEegSolar_csv); do \
		printf '\COPY "AnlagenEegSolar" FROM %s CSV;\n' "'$$f'" >>"$@"; \
	done
	echo 'COMMIT;' >>$@

AnlagenEegWind-csv.gen.sql: AnlagenEegWind.gen.sql $(AnlagenEegWind_csv)
	echo 'BEGIN;' >$@
	cat AnlagenEegWind.gen.sql >>$@
	for f in $(AnlagenEegWind_csv); do \
		printf '\COPY "AnlagenEegWind" FROM %s CSV;\n' "'$$f'" >>"$@"; \
	done
	echo 'COMMIT;' >>$@

AnlagenEegWasser-csv.gen.sql: AnlagenEegWasser.gen.sql $(AnlagenEegWasser_csv)
	echo 'BEGIN;' >$@
	cat AnlagenEegWasser.gen.sql >>$@
	for f in $(AnlagenEegWasser_csv); do \
		printf '\COPY "AnlagenEegWasser" FROM %s CSV;\n' "'$$f'" >>"$@"; \
	done
	echo 'COMMIT;' >>$@

AnlagenStromSpeicher-csv.gen.sql: AnlagenStromSpeicher.gen.sql $(AnlagenStromSpeicher_csv)
	echo 'BEGIN;' >$@
	cat AnlagenStromSpeicher.gen.sql >>$@
	for f in $(AnlagenStromSpeicher_csv); do \
		printf '\COPY "AnlagenStromSpeicher" FROM %s CSV;\n' "'$$f'" >>"$@"; \
	done
	echo 'COMMIT;' >>$@

AnlagenEegBiomasse-csv.gen.sql: AnlagenEegBiomasse.gen.sql $(AnlagenEegBiomasse_csv)
	echo 'BEGIN;' >$@
	cat AnlagenEegBiomasse.gen.sql >>$@
	for f in $(AnlagenEegBiomasse_csv); do \
		printf '\COPY "AnlagenEegBiomasse" FROM %s CSV;\n' "'$$f'" >>"$@"; \
	done
	echo 'COMMIT;' >>$@

Lokationen-csv.gen.sql: Lokationen.gen.sql $(Lokationen_csv)
	echo 'BEGIN;' >$@
	cat Lokationen.gen.sql >>$@
	for f in $(Lokationen_csv); do \
		printf '\COPY "Lokationen" FROM %s CSV;\n' "'$$f'" >>"$@"; \
	done
	echo 'COMMIT;' >>$@

Marktakteure-csv.gen.sql: Marktakteure.gen.sql $(Marktakteure_csv)
	echo 'BEGIN;' >$@
	cat Marktakteure.gen.sql >>$@
	for f in $(Marktakteure_csv); do \
		printf '\COPY "Marktakteure" FROM %s CSV;\n' "'$$f'" >>"$@"; \
	done
	echo 'COMMIT;' >>$@

Katalogkategorien-csv.gen.sql: Katalogkategorien.gen.sql $(Katalogkategorien_csv)
	echo 'BEGIN;' >$@
	cat Katalogkategorien.gen.sql >>$@
	for f in $(Katalogkategorien_csv); do \
		printf '\COPY "Katalogkategorien" FROM %s CSV;\n' "'$$f'" >>"$@"; \
	done
	echo 'COMMIT;' >>$@

Katalogwerte-csv.gen.sql: Katalogwerte.gen.sql $(Katalogwerte_csv)
	echo 'BEGIN;' >$@
	cat Katalogwerte.gen.sql >>$@
	for f in $(Katalogwerte_csv); do \
		printf '\COPY "Katalogwerte" FROM %s CSV;\n' "'$$f'" >>"$@"; \
	done
	echo 'COMMIT;' >>$@

Marktstammdatenregister.gen.sql: $(Marktstammdatenregister_sql)
	cat $(Marktstammdatenregister_sql) >"$@"

.PHONY: pg-insert
pg-insert: Marktstammdatenregister.gen.sql
	psql -h localhost -U user db <$<

Marktstammdatenregister.db: Marktstammdatenregister.gen.sql
	rm -f $@
	sqlite3 $@ <$<
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
	rm -f $(EinheitenBiomasse_xml) $(EinheitenBiomasse_csv)
	rm -f $(AnlagenEegSolar_xml) $(AnlagenEegSolar_csv)
	rm -f $(AnlagenEegWind_xml) $(AnlagenEegWind_csv)
	rm -f $(AnlagenEegWasser_xml) $(AnlagenEegWasser_csv)
	rm -f $(AnlagenStromSpeicher_xml) $(AnlagenStromSpeicher_csv)
	rm -f $(AnlagenEegBiomasse_xml) $(AnlagenEegBiomasse_csv)
	rm -f $(Lokationen_xml) $(Lokationen_csv)
	rm -f $(Katalogwerte_xml) $(Katalogwerte_csv)
	rm -f $(Katalogkategorien_xml) $(Katalogkategorien_csv)
	rm -f $(Marktakteure_xml) $(Marktakteure_csv)
	rm -f $(Marktstammdatenregister_sql)
	rm -f Marktstammdatenregister.db{,.br}
