# Building the database EinheitenSolar.db is a two-step process:
#
#     make download
#     make db

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

Marktakteure_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'Marktakteure_')
Marktakteure_csv := $(Marktakteure_xml:%.xml=%.csv)

Katalogkategorien_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'Katalogkategorien')
Katalogkategorien_csv := $(Katalogkategorien_xml:%.xml=%.csv)

Katalogwerte_xml := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'Katalogwerte')
Katalogwerte_csv := $(Katalogwerte_xml:%.xml=%.csv)

Marktstammdatenregister_sql := EinheitenSolar-csv.sql Marktakteure-csv.sql Katalogkategorien-csv.sql Katalogwerte-csv.sql

%.xsd: %.yml XSD.mustache
	mustache $< XSD.mustache >$@

%.sql: %.yml SQLite.mustache
	mustache $< SQLite.mustache >$@

%.jq: %.yml JQ.mustache
	mustache $< JQ.mustache >$@

# This rule unzips a single XML file, transforms it to JSON, and then
# transforms it to CSV.
#
# We perform all of this in a single rule so we can delete the XML and JSON
# files immediately. They are large and only necessary to create the (much
# smaller) CSV files.
EinheitenSolar_%.csv: Gesamtdatenexport.zip EinheitenSolar.xsd EinheitenSolar.jq
	unzip -qo $< $(subst ,,$(@:%.csv=%.xml))
	xmlschema-xml2json --schema=EinheitenSolar.xsd --force $(@:%.csv=%.xml)
	rm $(@:%.csv=%.xml)
	jq -r -f EinheitenSolar.jq $(@:%.csv=%.json) >$@
	rm $(@:%.csv=%.json)

Marktakteure_%.csv: Gesamtdatenexport.zip Marktakteure.xsd Marktakteure.jq
	unzip -qo $< $(subst ,,$(@:%.csv=%.xml))
	xmlschema-xml2json --schema=Marktakteure.xsd --force $(@:%.csv=%.xml)
	rm $(@:%.csv=%.xml)
	jq -r -f Marktakteure.jq $(@:%.csv=%.json) >$@
	rm $(@:%.csv=%.json)

Katalogkategorien.csv: Gesamtdatenexport.zip Katalogkategorien.xsd Katalogkategorien.jq
	unzip -qo $< $(subst ,,$(@:%.csv=%.xml))
	xmlschema-xml2json --schema=Katalogkategorien.xsd --force $(@:%.csv=%.xml)
	rm $(@:%.csv=%.xml)
	jq -r -f Katalogkategorien.jq $(@:%.csv=%.json) >$@
	rm $(@:%.csv=%.json)

Katalogwerte.csv: Gesamtdatenexport.zip Katalogwerte.xsd Katalogwerte.jq
	unzip -qo $< $(subst ,,$(@:%.csv=%.xml))
	xmlschema-xml2json --schema=Katalogwerte.xsd --force $(@:%.csv=%.xml)
	rm $(@:%.csv=%.xml)
	jq -r -f Katalogwerte.jq $(@:%.csv=%.json) >$@
	rm $(@:%.csv=%.json)

EinheitenSolar-csv.sql: EinheitenSolar.sql $(EinheitenSolar_csv)
	cp EinheitenSolar.sql $@
	echo '.mode csv' >>$@
	printf "%s\0" $(EinheitenSolar_csv) | \
		xargs -0 --no-run-if-empty --max-args=1 -I'{}' sh -c 'echo ".import {} EinheitenSolar --skip 1" >>$@'

Marktakteure-csv.sql: Marktakteure.sql $(Marktakteure_csv)
	cp Marktakteure.sql $@
	echo '.mode csv' >>$@
	printf "%s\0" $(Marktakteure_csv) | \
		xargs -0 --no-run-if-empty --max-args=1 -I'{}' sh -c 'echo ".import {} Marktakteure --skip 1" >>$@'

Katalogkategorien-csv.sql: Katalogkategorien.sql $(Katalogkategorien_csv)
	cp Katalogkategorien.sql $@
	echo '.mode csv' >>$@
	printf "%s\0" $(Katalogkategorien_csv) | \
		xargs -0 --no-run-if-empty --max-args=1 -I'{}' sh -c 'echo ".import {} Katalogkategorien --skip 1" >>$@'

Katalogwerte-csv.sql: Katalogwerte.sql $(Katalogwerte_csv)
	cp Katalogwerte.sql $@
	echo '.mode csv' >>$@
	printf "%s\0" $(Katalogwerte_csv) | \
		xargs -0 --no-run-if-empty --max-args=1 -I'{}' sh -c 'echo ".import {} Katalogwerte --skip 1" >>$@'

Marktstammdatenregister.db: $(Marktstammdatenregister_sql)
	rm -f $@
	cat $(Marktstammdatenregister_sql) | sqlite3 $@
	sqlite3 $@ <<<'CREATE INDEX idx_EinheitenSolar_Ort_Nettonennleistung on EinheitenSolar(Ort, Nettonennleistung);'
	sqlite3 $@ <<<'CREATE INDEX idx_EinheitenSolar_Inbetriebnahmedatum_Nettonennleistung on EinheitenSolar(Inbetriebnahmedatum, Nettonennleistung);'
	sqlite3 $@ <<<'CREATE INDEX idx_EinheitenSolar_DatumEndgueltigeStilllegung_Nettonennleistung on EinheitenSolar(DatumEndgueltigeStilllegung, Nettonennleistung);'
	sqlite3 $@ <<<'CREATE INDEX idx_EinheitenSolar_EinheitBetriebsstatus on EinheitenSolar(EinheitBetriebsstatus);'
	sqlite3 $@ <<<'CREATE INDEX idx_EinheitenSolar_Bundesland on EinheitenSolar(Bundesland);'
	sqlite3 $@ <<<'CREATE INDEX idx_EinheitenSolar_Landkreis on EinheitenSolar(Landkreis);'
	sqlite3 $@ <<<'CREATE INDEX idx_EinheitenSolar_Nutzungsbereich on EinheitenSolar(Nutzungsbereich);'
	sqlite3 $@ <<<'VACUUM; ANALYZE;'

Marktstammdatenregister.db.br: Marktstammdatenregister.db
	brotli -4 --keep --force --output=$@ $<
	touch $@

.PHONY: clean
clean:
	rm -f EinheitenSolar.{xsd,sql,jq}
	rm -f Marktakteure.{xsd,sql,jq}
	rm -f $(EinheitenSolar_xml) $(EinheitenSolar_csv)
	rm -f $(Marktakteure_xml) $(Marktakteure_csv)
	rm -f $(Marktstammdatenregister_sql)
	rm -f Marktstammdatenregister.db{,.br}
