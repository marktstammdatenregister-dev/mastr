#
# Build the database.
#
FROM ubuntu:20.04@sha256:1e48201ccc2ab83afc435394b3bf70af0fa0055215c1e26a5da9b50a1ae367c9 as builder
RUN apt-get -qq update \
 && apt-get -qq install \
      gdal-bin \
      spatialite-bin \
      wget \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /work
ARG OSM_URL
COPY ./osmconf.ini .
RUN wget --output-document raw.osm.pbf --no-verbose "${OSM_URL}" \
 && ogr2ogr -f SQLite osm.db raw.osm.pbf \
      -where "building is not null or boundary = 'administrative'" multipolygons \
      --config OSM_CONFIG_FILE osmconf.ini \
      -dsco SPATIALITE=YES \
      -gt 65536 \
 && rm raw.osm.pbf

COPY ./multipolygons-area.sql .
RUN spatialite osm.db <multipolygons-area.sql

#
# Run datasette.
#
FROM datasetteproject/datasette:0.58.1@sha256:e8749dd66c79c1808c37746469ecf73b816df515283745b4a5d53ce7f8f9c873 AS datasette
RUN apt-get -qq update \
 && apt-get -qq install dumb-init \
 && rm -rf /var/lib/apt/lists/*
RUN pip install \
      datasette-cluster-map \
      https://github.com/curiousleo/datasette-leaflet-geojson/archive/1e402abeb77192e0b8d51504b46055f1e1b4cf4d.tar.gz \
      datasette-vega \
 && true

WORKDIR /work
RUN groupadd -r datasette && useradd --no-log-init -r -g datasette datasette
RUN chown datasette:datasette .
USER datasette:datasette

COPY --from=builder /work/osm.db .
COPY ./metadata.yaml .
COPY ./settings.json .

ENTRYPOINT ["dumb-init", "--"]

# Use "." as the configuration directory. This loads *.db, metadata.yaml, settings.json, etc.
# https://docs.datasette.io/en/stable/settings.html#configuration-directory-mode
CMD ["datasette", "--port=8080", "--host=0.0.0.0", "--load-extension=spatialite", "."]
