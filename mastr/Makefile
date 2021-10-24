# Usage:
#
#     make download
#     make do-pgimport
#     make Marktstammdatenregister.dump

SHELL := bash

db_host := localhost
db_port := 5432
db_user := user
db_pass := pass
db_name := Marktstammdatenregister
db_url := 'postgres://$(db_user):$(db_pass)@$(db_host):$(db_port)/$(db_name)'
spec_dir := 'spec'

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

pgimport: cmd/pgimport/main.go
	go build -o $@ $<

%.xml: Gesamtdatenexport.zip
	unzip -qo "$<" "$@"
	mv "$@" "$@.utf16"
	iconv -f utf-16 -t utf-8 <"$@.utf16" >"$@"
	rm "$@.utf16"

.PHONY: do-pgimport
do-pgimport: Gesamtdatenexport.zip pgimport
	./pgimport -export $< -database $(db_url) -spec $(spec_dir)/Lokationen.yaml -prefix Lokationen
	./pgimport -export $< -database $(db_url) -spec $(spec_dir)/Katalogkategorien.yaml -prefix Katalogkategorien
	./pgimport -export $< -database $(db_url) -spec $(spec_dir)/Katalogwerte.yaml -prefix Katalogwerte
	./pgimport -export $< -database $(db_url) -spec $(spec_dir)/Marktakteure.yaml -prefix Marktakteure
	./pgimport -export $< -database $(db_url) -spec $(spec_dir)/AnlagenEegBiomasse.yaml -prefix AnlagenEegBiomasse
	./pgimport -export $< -database $(db_url) -spec $(spec_dir)/AnlagenEegSolar.yaml -prefix AnlagenEegSolar
	./pgimport -export $< -database $(db_url) -spec $(spec_dir)/AnlagenEegWasser.yaml -prefix AnlagenEegWasser
	./pgimport -export $< -database $(db_url) -spec $(spec_dir)/AnlagenEegWind.yaml -prefix AnlagenEegWind
	./pgimport -export $< -database $(db_url) -spec $(spec_dir)/AnlagenStromSpeicher.yaml -prefix AnlagenStromSpeicher
	./pgimport -export $< -database $(db_url) -spec $(spec_dir)/EinheitenBiomasse.yaml -prefix EinheitenBiomasse
	./pgimport -export $< -database $(db_url) -spec $(spec_dir)/EinheitenSolar.yaml -prefix EinheitenSolar
	./pgimport -export $< -database $(db_url) -spec $(spec_dir)/EinheitenStromSpeicher.yaml -prefix EinheitenStromSpeicher
	./pgimport -export $< -database $(db_url) -spec $(spec_dir)/EinheitenWasser.yaml -prefix EinheitenWasser
	./pgimport -export $< -database $(db_url) -spec $(spec_dir)/EinheitenWind.yaml -prefix EinheitenWind

Marktstammdatenregister.dump:
	PGPASSWORD=$(db_pass) pg_dump -h $(db_host) -U $(db_user) -Fc $(db_name) | pv -W >"$@"

.PHONY: restore-db
restore-db: Marktstammdatenregister.dump
	PGPASSWORD=$(db_pass) pg_restore -h $(db_host) -U $(db_user) -d $(db_name) "$<"

.PHONY: clean
clean:
	rm *.xml
	rm -f Marktstammdatenregister.dump
