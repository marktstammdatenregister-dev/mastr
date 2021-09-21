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

Marktstammdatenregister := (EinheitenSolar-csv.sql Marktakteure_csv.sql)

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
	rm -f $(@:%.csv=%.json)
	xmlschema-xml2json --schema=EinheitenSolar.xsd $(@:%.csv=%.xml)
	rm $(@:%.csv=%.xml)
	jq -r -f EinheitenSolar.jq $(@:%.csv=%.json) >$@
	rm $(@:%.csv=%.json)

Marktakteure_%.csv: Gesamtdatenexport.zip Marktakteure.xsd Marktakteure.jq
	unzip -qo $< $(subst ,,$(@:%.csv=%.xml))
	rm -f $(@:%.csv=%.json)
	xmlschema-xml2json --schema=Marktakteure.xsd $(@:%.csv=%.xml)
	rm $(@:%.csv=%.xml)
	jq -r -f Marktakteure.jq $(@:%.csv=%.json) >$@
	rm $(@:%.csv=%.json)

EinheitenSolar-csv.sql: EinheitenSolar.sql $(EinheitenSolar_csv)
	cp EinheitenSolar.sql $@
	echo '.mode csv' >>$@
	printf "%s\0" $(EinheitenSolar_csv) | \
		xargs -0 --no-run-if-empty --max-args=1 -I'{}' sh -c 'echo ".import {} EinheitenSolar --skip 1" >>$@'

Marktakteure-csv.sql: Marktakteure.sql $(Marktakteure_csv)
	cp Marktakteure.sql $@
	printf "%s\0" $(Marktakteure_csv) | \
		xargs -0 --no-run-if-empty --max-args=1 -I'{}' sh -c 'echo ".import {} Marktakteure --skip 1" >>$@'

Marktstammdatenregister.db: $(Marktstammdatenregister)
	rm -f $@
	cat $(Marktstammdatenregister) | sqlite3 $@
	sqlite3 $@ <<<'VACUUM; ANALYZE;'

Marktstammdatenregister.db.br: Marktstammdatenregister.db
	brotli -4 --keep --force --output=$@ $<
	touch $@
