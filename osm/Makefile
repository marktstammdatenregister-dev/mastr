# Download input.osm.pbf from https://download.geofabrik.de/ before running 'make'.

buildings.db: germany-latest.osm.pbf buildings.cfg buildings.sql
	osmium tags-filter germany-latest.osm.pbf a/building --output buildings.osm.opl --remove-tags --overwrite
	osmium export buildings.osm.opl --config buildings.cfg --output buildings-broken-ewkb.pg --overwrite
	pv buildings-broken-ewkb.pg | ./fix_ewkb.py >buildings.pg
	spatialite $@ <buildings.sql

boundaries.db: germany-latest.osm.pbf boundaries.cfg boundaries.sql
	osmium tags-filter germany-latest.osm.pbf a/boundary=administrative --output boundaries.osm.opl --remove-tags --overwrite
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
