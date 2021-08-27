#
# Build the database.
#
FROM ubuntu:20.04@sha256:1e48201ccc2ab83afc435394b3bf70af0fa0055215c1e26a5da9b50a1ae367c9 as builder
RUN apt-get -qq update \
 && apt-get -qq install \
      gdal-bin \
      spatialite-bin \
      wget \
 && true

WORKDIR /work
ARG OSM_URL
RUN wget --output-document raw.osm.pbf --no-verbose "${OSM_URL}"
RUN ogr2ogr -f SQLite osm.db raw.osm.pbf -dsco SPATIALITE=YES -gt 65536 -where "building is not null or boundary = 'administrative'" multipolygons

COPY ./multipolygons-area.sql .
RUN spatialite osm.db <multipolygons-area.sql

#
# Run datasette.
#
FROM datasetteproject/datasette:0.58.1@sha256:e8749dd66c79c1808c37746469ecf73b816df515283745b4a5d53ce7f8f9c873 AS datasette
RUN apt-get -qq update \
 && apt-get -qq install dumb-init
RUN pip install datasette-leaflet-geojson

WORKDIR /work
RUN groupadd -r datasette && useradd --no-log-init -r -g datasette datasette
RUN chown datasette:datasette .
USER datasette:datasette

COPY --from=builder /work/osm.db .
COPY ./metadata.yaml .

ENTRYPOINT ["dumb-init", "--"]
CMD ["datasette", "--port=8080", "--host=0.0.0.0", "/work/osm.db", "--load-extension=spatialite", "--metadata=/work/metadata.yaml"]
