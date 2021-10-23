# Usage:
#
#     make download
#     make populate-db
#     make Marktstammdatenregister.dump

SHELL := bash

db_host := localhost
db_port := 5432
db_user := user
db_pass := pass
db_name := Marktstammdatenregister
db_url := 'postgres://$(db_user):$(db_pass)@$(db_host):$(db_port)/$(db_name)'
#db_url := 'postgres://$(db_name)?host=/tmp/postgresql/run&user=$(db_user)&database=$(db_name)'

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
$(Lokationen_xml) \
$(Katalogkategorien_xml) \
$(Katalogwerte_xml) \
$(Marktakteure_xml) \
$(AnlagenEegBiomasse_xml) \
$(AnlagenEegSolar_xml) \
$(AnlagenEegWasser_xml) \
$(AnlagenEegWind_xml) \
$(AnlagenStromSpeicher_xml) \
$(EinheitenBiomasse_xml) \
$(EinheitenSolar_xml) \
$(EinheitenStromSpeicher_xml) \
$(EinheitenWasser_xml) \
$(EinheitenWind_xml)

xml2psql: xml2psql.go
	go build -o $@ $<

%.xml: Gesamtdatenexport.zip
	unzip -qo "$<" "$@"
	mv "$@" "$@.utf16"
	iconv -f utf-16 -t utf-8 <"$@.utf16" >"$@"
	rm "$@.utf16"

.PHONY: populate-db
populate-db: Gesamtdatenexport.zip xml2psql
	./xml2psql -export $< -database $(db_url) -descriptor Lokationen.yaml -prefix Lokationen
	./xml2psql -export $< -database $(db_url) -descriptor Katalogkategorien.yaml -prefix Katalogkategorien
	./xml2psql -export $< -database $(db_url) -descriptor Katalogwerte.yaml -prefix Katalogwerte
	./xml2psql -export $< -database $(db_url) -descriptor Marktakteure.yaml -prefix Marktakteure
	./xml2psql -export $< -database $(db_url) -descriptor AnlagenEegBiomasse.yaml -prefix AnlagenEegBiomasse
	./xml2psql -export $< -database $(db_url) -descriptor AnlagenEegSolar.yaml -prefix AnlagenEegSolar
	./xml2psql -export $< -database $(db_url) -descriptor AnlagenEegWasser.yaml -prefix AnlagenEegWasser
	./xml2psql -export $< -database $(db_url) -descriptor AnlagenEegWind.yaml -prefix AnlagenEegWind
	./xml2psql -export $< -database $(db_url) -descriptor AnlagenStromSpeicher.yaml -prefix AnlagenStromSpeicher
	./xml2psql -export $< -database $(db_url) -descriptor EinheitenBiomasse.yaml -prefix EinheitenBiomasse
	./xml2psql -export $< -database $(db_url) -descriptor EinheitenSolar.yaml -prefix EinheitenSolar
	./xml2psql -export $< -database $(db_url) -descriptor EinheitenStromSpeicher.yaml -prefix EinheitenStromSpeicher
	./xml2psql -export $< -database $(db_url) -descriptor EinheitenWasser.yaml -prefix EinheitenWasser
	./xml2psql -export $< -database $(db_url) -descriptor EinheitenWind.yaml -prefix EinheitenWind

.PHONY: db-Lokationen db-Katalogkategorien db-Katalogwerte db-Marktakteure db-AnlagenEegBiomasse db-AnlagenEegSolar db-AnlagenEegWind db-AnlagenStromSpeicher db-EinheitenBiomasse db-EinheitenSolar db-EinheitenStromSpeicher db-EinheitenWasser db-EinheitenWind
db-Lokationen: Gesamtdatenexport.zip xml2psql
	./xml2psql -export $< -database $(db_url) -descriptor Lokationen.yaml -prefix Lokationen

db-Katalogkategorien: Gesamtdatenexport.zip xml2psql
	./xml2psql -export $< -database $(db_url) -descriptor Katalogkategorien.yaml -prefix Katalogkategorien

db-Katalogwerte: Gesamtdatenexport.zip xml2psql db-Katalogkategorien
	./xml2psql -export $< -database $(db_url) -descriptor Katalogwerte.yaml -prefix Katalogwerte

db-Marktakteure: Gesamtdatenexport.zip xml2psql db-Katalogwerte
	./xml2psql -export $< -database $(db_url) -descriptor Marktakteure.yaml -prefix Marktakteure

db-AnlagenEegBiomasse: Gesamtdatenexport.zip xml2psql db-Katalogwerte
	./xml2psql -export $< -database $(db_url) -descriptor AnlagenEegBiomasse.yaml -prefix AnlagenEegBiomasse

db-AnlagenEegSolar: Gesamtdatenexport.zip xml2psql db-Katalogwerte
	./xml2psql -export $< -database $(db_url) -descriptor AnlagenEegSolar.yaml -prefix AnlagenEegSolar

db-AnlagenEegWasser: Gesamtdatenexport.zip xml2psql db-Katalogwerte
	./xml2psql -export $< -database $(db_url) -descriptor AnlagenEegWasser.yaml -prefix AnlagenEegWasser

db-AnlagenEegWind: Gesamtdatenexport.zip xml2psql db-Katalogwerte
	./xml2psql -export $< -database $(db_url) -descriptor AnlagenEegWind.yaml -prefix AnlagenEegWind

db-AnlagenStromSpeicher: Gesamtdatenexport.zip xml2psql db-Katalogwerte
	./xml2psql -export $< -database $(db_url) -descriptor AnlagenStromSpeicher.yaml -prefix AnlagenStromSpeicher

db-EinheitenBiomasse: Gesamtdatenexport.zip xml2psql db-Lokationen db-Marktakteure db-Katalogwerte db-AnlagenEegBiomasse
	./xml2psql -export $< -database $(db_url) -descriptor EinheitenBiomasse.yaml -prefix EinheitenBiomasse

db-EinheitenSolar: Gesamtdatenexport.zip xml2psql db-Lokationen db-Marktakteure db-Katalogwerte db-AnlagenEegSolar
	./xml2psql -export $< -database $(db_url) -descriptor EinheitenSolar.yaml -prefix EinheitenSolar

db-EinheitenStromSpeicher: Gesamtdatenexport.zip xml2psql db-Lokationen db-Marktakteure db-Katalogwerte db-AnlagenStromSpeicher
	./xml2psql -export $< -database $(db_url) -descriptor EinheitenStromSpeicher.yaml -prefix EinheitenStromSpeicher

db-EinheitenWasser: Gesamtdatenexport.zip xml2psql db-Lokationen db-Marktakteure db-Katalogwerte db-AnlagenEegWasser
	./xml2psql -export $< -database $(db_url) -descriptor EinheitenWasser.yaml -prefix EinheitenWasser

db-EinheitenWind: Gesamtdatenexport.zip xml2psql db-Lokationen db-Marktakteure db-Katalogwerte db-AnlagenEegWind
	./xml2psql -export $< -database $(db_url) -descriptor EinheitenWind.yaml -prefix EinheitenWind

.PHONY: populate-db2
populate-db2: db-Lokationen db-Katalogkategorien db-Katalogwerte db-Marktakteure db-AnlagenEegBiomasse db-AnlagenEegSolar db-AnlagenEegWind db-AnlagenStromSpeicher db-EinheitenBiomasse db-EinheitenSolar db-EinheitenStromSpeicher db-EinheitenWasser db-EinheitenWind

Marktstammdatenregister.dump:
	PGPASSWORD=$(db_pass) pg_dump -h $(db_host) -U $(db_user) -Fc $(db_name) | pv -W >"$@"

.PHONY: restore-db
restore-db: Marktstammdatenregister.dump
	PGPASSWORD=$(db_pass) pg_restore -h $(db_host) -U $(db_user) -d $(db_name) "$<"

.PHONY: clean
clean:
	rm *.xml
	rm -f Marktstammdatenregister.dump
