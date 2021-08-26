FROM ubuntu:20.04@sha256:1e48201ccc2ab83afc435394b3bf70af0fa0055215c1e26a5da9b50a1ae367c9 as builder
RUN apt-get -qq update \
 && apt-get -qq install \
      gdal-bin \
      spatialite-bin \
      wget \
 && true

WORKDIR /work
RUN wget -O raw.osm.pbf https://download.geofabrik.de/europe/germany/baden-wuerttemberg-210824.osm.pbf
RUN ogr2ogr -f SQLite raw.db raw.osm.pbf -dsco SPATIALITE=YES -gt 65536

COPY ./preprocess.sql .
RUN ls -lh raw.db
RUN spatialite raw.db <preprocess.sql
RUN ls -lh raw.db


FROM datasetteproject/datasette:0.58.1@sha256:e8749dd66c79c1808c37746469ecf73b816df515283745b4a5d53ce7f8f9c873 AS datasette

WORKDIR /work
COPY --from=builder /work/raw.db osm.db
COPY ./metadata.yaml .

RUN pip install datasette-leaflet-geojson
