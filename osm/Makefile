# Download input.osm.pbf from https://download.geofabrik.de/ before running 'make'.

ewkb2wkt: ewkb2wkt.go
	go build -o "$@" "$<"

buildings.osm.pbf: germany-latest.osm.pbf
	osmium tags-filter "$<" a/building --output "$@" --remove-tags --overwrite

buildings.tsv: buildings.osm.pbf ewkb2wkt
	osmium export "$<" --output-format pg | cut -f1 | ./ewkb2wkt -min-area=1000 | pv >"$@"

boundaries.osm.pbf: germany-latest.osm.pbf
	osmium tags-filter "$<" a/boundary=administrative --output "$@" --remove-tags --overwrite

boundaries.tsv: boundaries.osm.pbf ewkb2wkt
	osmium export "$<" --output-format pg | ./ewkb2wkt | pv >"$@"

OpenStreetMap.db: buildings.tsv buildings.sql boundaries.tsv boundaries.sql
	rm -f "$@"
	spatialite "$@" <buildings.sql
	spatialite "$@" <boundaries.sql
	spatialite "$@" 'analyze; vacuum;'

%.br: %
	brotli -4 --keep --force --output="$@" "$<"

.PHONY: all
all: OpenStreetMap.db.br

.PHONY: clean-intermediate
clean-intermediate:
	rm -f *.tsv {boundaries,buildings}.osm.pbf
