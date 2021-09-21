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

#
# Run datasette.
#
FROM docker.io/datasetteproject/datasette:0.58.1@sha256:e8749dd66c79c1808c37746469ecf73b816df515283745b4a5d53ce7f8f9c873 AS datasette
RUN apt-get -qq update \
 && apt-get -qq install \
      brotli \
      dumb-init \
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

#COPY --from=builder-osm /work/boundaries.db.br .
#COPY --from=builder-osm /work/buildings.db.br .
#COPY --from=builder-mastr /work/Marktstammdatenregister.db.br .
COPY ./mastr/Marktstammdatenregister.db.br .
COPY ./metadata.yaml .
COPY ./settings.json .

ENTRYPOINT ["dumb-init", "--"]

# Use "." as the configuration directory. This loads *.db, metadata.yaml, settings.json, etc.
# https://docs.datasette.io/en/stable/settings.html#configuration-directory-mode
CMD ["sh", "-c", "brotli --rm --decompress --no-copy-stat *.db.br && datasette --port=8080 --host=0.0.0.0 --load-extension=spatialite --cors ."]
