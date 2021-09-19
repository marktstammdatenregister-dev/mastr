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

# Run 'make download' to populate 'xml_files' and 'csv_files'.
xml_files := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'EinheitenSolar_')
csv_files := $(xml_files:%.xml=%.csv)

# This rule unzips a single XML file, transforms it to JSON, and then
# transforms it to CSV.
#
# We perform all of this in a single rule so we can delete the XML and JSON
# files immediately. They are large and only necessary to create the (much
# smaller) CSV files.
%.csv: Gesamtdatenexport.zip EinheitenSolar.xsd EinheitenSolar.jq
	unzip -qo $< $(subst ,,$(@:%.csv=%.xml))
	rm -f $(@:%.csv=%.json)
	xmlschema-xml2json --schema=EinheitenSolar.xsd $(@:%.csv=%.xml)
	rm $(@:%.csv=%.xml)
	jq -r -f EinheitenSolar.jq $(@:%.csv=%.json) >$@
	rm $(@:%.csv=%.json)

EinheitenSolar-csv.sql: EinheitenSolar.sql $(csv_files)
	cp EinheitenSolar.sql $@
	echo '.mode csv' >>$@
	printf "%s\0" $(csv_files) | \
		xargs -0 --no-run-if-empty --max-args=1 -I'{}' sh -c 'echo ".import {} EinheitenSolar --skip 1" >>$@'

EinheitenSolar-raw.db: EinheitenSolar-csv.sql
	rm -f $@
	sqlite3 $@ <$<

EinheitenSolar.db: EinheitenSolar-raw.db EinheitenSolar-points.sql
	cp EinheitenSolar{-raw,}.db
	spatialite $@ <EinheitenSolar-points.sql
	sqlite3 $@ <<<'VACUUM; ANALYZE;'

EinheitenSolar.db.br: EinheitenSolar.db
	brotli -4 --keep --force --output=$@ $<
	touch $@

.PHONY: db
db: EinheitenSolar.db

.PHONY: image
image: EinheitenSolar.db.br
	docker build . -t mastr
