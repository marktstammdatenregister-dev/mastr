# Download input.osm.pbf from https://download.geofabrik.de/ before running 'make'.

ewkb2wkt: ewkb2wkt.go
	go build -o "$@" "$<"

buildings.tsv: germany-latest.osm.pbf ewkb2wkt
	osmium tags-filter germany-latest.osm.pbf a/building type=multipolygon --output buildings.osm.pbf --remove-tags --overwrite
	osmium export buildings.osm.pbf --output-format pg | cut -f1 | ./ewkb2wkt --min-area=1000 | pv >"$@"
	rm buildings.osm.pbf

boundaries.tsv: germany-latest.osm.pbf ewkb2wkt
	osmium tags-filter germany-latest.osm.pbf a/boundary=administrative type=boundary --output boundaries.osm.pbf --remove-tags --overwrite
	osmium export boundaries.osm.pbf --output-format pg | ./ewkb2wkt | pv >"$@"
	rm boundaries.osm.pbf

OpenStreetMap.db: buildings.tsv buildings.sql boundaries.tsv boundaries.sql
	rm -f "$@"
	spatialite "$@" <buildings.sql
	spatialite "$@" <boundaries.sql

%.br: %
	brotli -4 --keep --force --output="$@" "$<"

.PHONY: all
all: OpenStreetMap.db.br

.PHONY: clean-intermediate
clean-intermediate:
	rm -f *.tsv *.osm.pbf
