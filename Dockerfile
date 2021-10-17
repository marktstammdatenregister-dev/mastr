##
## Build the OSM database.
##
#FROM docker.io/ubuntu:20.04@sha256:1e48201ccc2ab83afc435394b3bf70af0fa0055215c1e26a5da9b50a1ae367c9 as builder-osm
#RUN apt-get -qq update \
# && DEBIAN_FRONTEND=noninteractive apt-get -qq install --no-install-recommends \
#      brotli \
#      gdal-bin \
#      make \
#      spatialite-bin \
#      wget \
# && rm -rf /var/lib/apt/lists/*
#
#WORKDIR /work
#ARG OSM_URL
#RUN wget --output-document input.osm.pbf --no-verbose "${OSM_URL}"
#COPY ./osm/ ./
#RUN make -j all && make clean-intermediate
#
##
## Build the MaStR database.
##
#FROM docker.io/ubuntu:20.04@sha256:1e48201ccc2ab83afc435394b3bf70af0fa0055215c1e26a5da9b50a1ae367c9 as builder-mastr
#RUN apt-get -qq update \
# && DEBIAN_FRONTEND=noninteractive apt-get -qq install --no-install-recommends \
#      axel \
#      brotli \
#      curl \
#      jq \
#      make \
#      python3-pip \
#      spatialite-bin \
#      unzip \
#      wget \
# && rm -rf /var/lib/apt/lists/*
#RUN curl -sSL -o pup.zip https://github.com/ericchiang/pup/releases/download/v0.4.0/pup_v0.4.0_linux_amd64.zip \
# && unzip pup.zip \
# && rm pup.zip \
# && mv pup /usr/bin/pup \
# && chmod +x /usr/bin/pup
#RUN pip3 install xmlschema
#
#WORKDIR /work
#ARG OSM_URL
#COPY ./mastr/ ./
#RUN make download && make -j8 Marktstammdatenregister.db.br

# https://github.com/simonw/datasette/blob/63886178a649586b403966a27a45881709d2b868/Dockerfile
# But with bullseye instead of buster so we get a newer sqlite and spatialite version
# https://packages.debian.org/search?suite=bullseye&searchon=names&keywords=spatialite
FROM python:3.9-slim-bullseye as datasette

ARG VERSION=0.58.1

RUN apt-get update && \
    apt-get install -y --no-install-recommends libsqlite3-mod-spatialite

RUN pip install https://github.com/simonw/datasette/archive/refs/tags/${VERSION}.zip && \
    find /usr/local/lib -name '__pycache__' | xargs rm -r && \
    rm -rf /root/.cache/pip

EXPOSE 8001
CMD ["datasette"]

#
# Run datasette.
#
#FROM docker.io/datasetteproject/datasette:0.58.1@sha256:e8749dd66c79c1808c37746469ecf73b816df515283745b4a5d53ce7f8f9c873 AS datasette
FROM datasette as runner
RUN apt-get -qq update \
 && apt-get -qq install \
      brotli \
      dumb-init \
 && rm -rf /var/lib/apt/lists/*
# https://github.com/curiousleo/datasette-leaflet-geojson/archive/1e402abeb77192e0b8d51504b46055f1e1b4cf4d.tar.gz
RUN pip install \
      datasette-cluster-map \
      datasette-leaflet-geojson \
      datasette-vega \
 && true

WORKDIR /work
RUN groupadd -r datasette && useradd --no-log-init -r -g datasette datasette
RUN chown datasette:datasette .
USER datasette:datasette

FROM runner as copied

# The database files are renamed to .sqlite3 on purpose.
#
# The only way I can find in the documentation to load settings from a file is to use "Configuration
# directory mode", i.e. passing a directory to `datasette`.
#
# https://docs.datasette.io/en/0.58.1/settings.html#id2
#
# However, we also want users to be able to download the SQLite files. This requires us to pass
# `--immutable <dbfile>` to `datasette`.
#
# If we use `--immutable <dbfile>` in combination with "Configuration directory mode", Datasette
# will pick up any *.db files in the given directory and show them. If any of those SQLite files
# were passed in with `--immutable`, they will be shown twice: once as an immutable database, once
# as a mutable database.
#
# To avoid Datasette listing the databases twice, we give the SQLite files a file ending other than
# "db" so they are not picked up and listed as databases as a result of using "Configuration
# directory mode".
COPY ./mastr/Marktstammdatenregister.db.br ./Marktstammdatenregister.sqlite3.br
COPY ./osm/OpenStreetMap.db.br ./OpenStreetMap.sqlite3.br
COPY ./metadata.yaml .
COPY ./settings.json .

ENTRYPOINT ["dumb-init", "--"]

# Use "." as the configuration directory. This loads *.db, metadata.yaml, settings.json, etc.
# https://docs.datasette.io/en/stable/settings.html#configuration-directory-mode
CMD ["sh", "-c", "brotli --rm --decompress --no-copy-stat *.br && datasette --port=8080 --host=0.0.0.0 --load-extension=spatialite --cors --immutable Marktstammdatenregister.sqlite3 --immutable OpenStreetMap.sqlite3 ."]
