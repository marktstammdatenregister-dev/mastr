# Download input.osm.pbf from https://download.geofabrik.de/ before running 'make'.

buildings.db: input.osm.pbf buildings.ini multipolygons-area.sql
	rm -f "$@"
	ogr2ogr -f SQLite "$@" input.osm.pbf multipolygons \
		-where "building is not null" \
		--config OSM_CONFIG_FILE buildings.ini \
		-dsco SPATIALITE=YES \
		-gt 65536
	sqlite3 "$@" <multipolygons-area.sql

boundaries.db: input.osm.pbf boundaries.cfg boundaries.sql
	osmium tags-filter input.osm.pbf a/boundary=administrative --output boundaries.osm.opl --remove-tags --overwrite
	osmium export boundaries.osm.opl --config boundaries.cfg --output-format pg \
		| ./fix_ewkb.py >boundaries.pg
	spatialite $@ <boundaries.sql

%.br: %
	brotli -4 --keep --output="$@" "$<"

.PHONY: all
all: buildings.db.br boundaries.db.br

.PHONY: clean-intermediate
clean-intermediate:
	rm -f buildings.db boundaries.db
