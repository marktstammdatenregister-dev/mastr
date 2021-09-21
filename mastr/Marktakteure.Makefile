# Building the database Marktakteure.db is a two-step process:
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
xml_files := $(shell unzip -Z1 Gesamtdatenexport.zip | grep 'Marktakteure_')
csv_files := $(xml_files:%.xml=%.csv)

# This rule unzips a single XML file, transforms it to JSON, and then
# transforms it to CSV.
#
# We perform all of this in a single rule so we can delete the XML and JSON
# files immediately. They are large and only necessary to create the (much
# smaller) CSV files.
Marktakteure_%.csv: Gesamtdatenexport.zip Marktakteure.xsd Marktakteure.jq
	unzip -qo $< $(subst ,,$(@:%.csv=%.xml))
	rm -f $(@:%.csv=%.json)
	xmlschema-xml2json --schema=Marktakteure.xsd $(@:%.csv=%.xml)
	rm $(@:%.csv=%.xml)
	jq -r -f Marktakteure.jq $(@:%.csv=%.json) >$@
	rm $(@:%.csv=%.json)

Marktakteure-csv.sql: Marktakteure.sql $(csv_files)
	cp Marktakteure.sql $@
	printf "%s\0" $(csv_files) | \
		xargs -0 --no-run-if-empty --max-args=1 -I'{}' sh -c 'echo ".import {} Marktakteure --skip 1" >>$@'

Marktakteure.db: Marktakteure-csv.sql
	rm -f $@
	sqlite3 $@ <$<
	sqlite3 $@ <<<'VACUUM; ANALYZE;'

Marktakteure.db.br: Marktakteure.db
	brotli -4 --keep --force --output=$@ $<
	touch $@

.PHONY: db
db: Marktakteure.db

.PHONY: image
image: Marktakteure.db.br
	docker build . -t mastr
