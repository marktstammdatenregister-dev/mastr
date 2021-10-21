# Usage:
#
#     make download
#     make populate-db

SHELL := bash

db_host := localhost
db_port := 5432
db_user := user
db_pass := pass
db_name := Marktstammdatenregister
db_url := 'postgres://$(db_user):$(db_pass)@$(db_host):$(db_port)/$(db_name)'

# Download the latest export.
Gesamtdatenexport.zip:
	curl -sSL https://www.marktstammdatenregister.de/MaStR/Datendownload | \
		pup 'a[href^="https://download.marktstammdatenregister.de/Gesamtdatenexport"][href$$=".zip"] attr{href}' | \
		xargs axel --quiet --output=$@

.PHONY: download
download: Gesamtdatenexport.zip

# Run 'make download' to populate these variables.
EinheitenSolar_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'EinheitenSolar_')
EinheitenWind_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'EinheitenWind')
EinheitenWasser_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'EinheitenWasser')
EinheitenStromSpeicher_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'EinheitenStromSpeicher_')
EinheitenBiomasse_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'EinheitenBiomasse')
AnlagenEegSolar_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'AnlagenEegSolar_')
AnlagenEegWind_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'AnlagenEegWind')
AnlagenEegWasser_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'AnlagenEegWasser')
AnlagenStromSpeicher_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'AnlagenStromSpeicher_')
AnlagenEegBiomasse_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'AnlagenEegBiomasse')
Lokationen_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'Lokationen_')
Marktakteure_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'Marktakteure_')
Katalogkategorien_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'Katalogkategorien')
Katalogwerte_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'Katalogwerte')

Marktstammdatenregister_xml := \
$(Lokationen) \
$(Katalogkategorien) \
$(Katalogwerte) \
$(Marktakteure) \
$(AnlagenEegBiomasse) \
$(AnlagenEegSolar) \
$(AnlagenEegWasser) \
$(AnlagenEegWind) \
$(AnlagenStromSpeicher) \
$(EinheitenBiomasse) \
$(EinheitenSolar) \
$(EinheitenStromSpeicher) \
$(EinheitenWasser) \
$(EinheitenWind)

xml2csv: xml2csv.go
	go build -o $@ $<

.PHONY: populate-db
populate-db: xml2csv $(Marktstammdatenregister_xml)
	find . -name 'Lokationen*.xml' | \
		xargs -I{} sh -c 'iconv -f utf-16 -t utf-8 {} | ./xml2csv -descriptor Lokationen.yaml -database $(db_url)'
	find . -name 'Katalogkategorien*.xml' | \
		xargs -I{} sh -c 'iconv -f utf-16 -t utf-8 {} | ./xml2csv -descriptor Katalogkategorien.yaml -database $(db_url)'
	find . -name 'Katalogwerte*.xml' | \
		xargs -I{} sh -c 'iconv -f utf-16 -t utf-8 {} | ./xml2csv -descriptor Katalogwerte.yaml -database $(db_url)'
	find . -name 'Marktakteure*.xml' | \
		xargs -I{} sh -c 'iconv -f utf-16 -t utf-8 {} | ./xml2csv -descriptor Marktakteure.yaml -database $(db_url)'
	find . -name 'AnlagenEegBiomasse*.xml' | \
		xargs -I{} sh -c 'iconv -f utf-16 -t utf-8 {} | ./xml2csv -descriptor AnlagenEegBiomasse.yaml -database $(db_url)'
	find . -name 'AnlagenEegSolar*.xml' | \
		xargs -I{} sh -c 'iconv -f utf-16 -t utf-8 {} | ./xml2csv -descriptor AnlagenEegSolar.yaml -database $(db_url)'
	find . -name 'AnlagenEegWasser*.xml' | \
		xargs -I{} sh -c 'iconv -f utf-16 -t utf-8 {} | ./xml2csv -descriptor AnlagenEegWasser.yaml -database $(db_url)'
	find . -name 'AnlagenEegWind*.xml' | \
		xargs -I{} sh -c 'iconv -f utf-16 -t utf-8 {} | ./xml2csv -descriptor AnlagenEegWind.yaml -database $(db_url)'
	find . -name 'AnlagenStromSpeicher*.xml' | \
		xargs -I{} sh -c 'iconv -f utf-16 -t utf-8 {} | ./xml2csv -descriptor AnlagenStromSpeicher.yaml -database $(db_url)'
	find . -name 'EinheitenBiomasse*.xml' | \
		xargs -I{} sh -c 'iconv -f utf-16 -t utf-8 {} | ./xml2csv -descriptor EinheitenBiomasse.yaml -database $(db_url)'
	find . -name 'EinheitenSolar*.xml' | \
		xargs -I{} sh -c 'iconv -f utf-16 -t utf-8 {} | ./xml2csv -descriptor EinheitenSolar.yaml -database $(db_url)'
	find . -name 'EinheitenStromSpeicher*.xml' | \
		xargs -I{} sh -c 'iconv -f utf-16 -t utf-8 {} | ./xml2csv -descriptor EinheitenStromSpeicher.yaml -database $(db_url)'
	find . -name 'EinheitenWasser*.xml' | \
		xargs -I{} sh -c 'iconv -f utf-16 -t utf-8 {} | ./xml2csv -descriptor EinheitenWasser.yaml -database $(db_url)'
	find . -name 'EinheitenWind*.xml' | \
		xargs -I{} sh -c 'iconv -f utf-16 -t utf-8 {} | ./xml2csv -descriptor EinheitenWind.yaml -database $(db_url)'

Marktstammdatenregister.dump:
	PGPASSWORD=$(db_pass) pg_dump -h $(db_host) -U $(db_user) -Fc $(db_name) | pv -W >"$@"

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

Marktstammdatenregister.dump.br: Marktstammdatenregister.dump
	brotli -4 --keep --force --output=$@ $<
	touch $@

.PHONY: clean
clean:
	rm -f Marktstammdatenregister.dump{,.br}
