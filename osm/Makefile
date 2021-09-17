# Download input.osm.pbf from https://download.geofabrik.de/ before running 'make'.

buildings.db: input.osm.pbf buildings.ini multipolygons-area.sql
	rm -f "$@"
	ogr2ogr -f SQLite "$@" input.osm.pbf multipolygons \
		-where "building is not null" \
		--config OSM_CONFIG_FILE buildings.ini \
		-dsco SPATIALITE=YES \
		-gt 65536
	spatialite "$@" <multipolygons-area.sql

boundaries.db: input.osm.pbf buildings.ini
	rm -f "$@"
	ogr2ogr -f SQLite "$@" input.osm.pbf multipolygons \
		-where "boundary = 'administrative'" \
		--config OSM_CONFIG_FILE buildings.ini \
		-dsco SPATIALITE=YES \
		-gt 65536

points.db: input.osm.pbf points.ini
	rm -f "$@"
	ogr2ogr -f SQLite "$@" input.osm.pbf points \
		-where "generator_source = 'solar'" \
		--config OSM_CONFIG_FILE points.ini \
		-dsco SPATIALITE=YES \
		-gt 65536

%.br: %
	brotli -4 --keep --output="$@" "$<"

.PHONY: all
all: buildings.db.br boundaries.db.br points.db.br

.PHONY: clean-intermediate
clean-intermediate:
	rm -f buildings.db boundaries.db points.db
