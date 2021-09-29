# Download input.osm.pbf from https://download.geofabrik.de/ before running 'make'.

ewkb2wkt: ewkb2wkt.go
	go build -o "$@" "$<"

buildings.tsv: germany-latest.osm.pbf buildings.cfg ewkb2wkt
	osmium tags-filter germany-latest.osm.pbf a/building --output buildings.osm.opl --remove-tags --overwrite
	osmium export buildings.osm.opl --config buildings.cfg --output-format pg \
		| cut -f1 | ./ewkb2wkt | pv >"$@"
	rm buildings.osm.opl

boundaries.tsv: germany-latest.osm.pbf boundaries.cfg ewkb2wkt
	osmium tags-filter germany-latest.osm.pbf a/boundary=administrative --output boundaries.osm.opl --remove-tags --overwrite
	osmium export boundaries.osm.opl --config boundaries.cfg --output-format pg \
		| ./ewkb2wkt | pv >"$@"
	rm boundaries.osm.opl

OpenStreetMap.db: buildings.tsv buildings.sql boundaries.tsv boundaries.sql
	spatialite "$@" <buildings.sql
	spatialite "$@" <boundaries.sql

%.br: %
	brotli -4 --keep --output="$@" "$<"

.PHONY: all
all: OpenStreetMap.db.br

.PHONY: clean-intermediate
clean-intermediate:
	rm -f *.tsv *.osm.opl
