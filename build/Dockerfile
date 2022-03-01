#
# Run datasette.
#
FROM docker.io/datasetteproject/datasette:0.60.1@sha256:c2f93ba54eb49e299dd98180de0a8ca3cc1da0df05425a8e32f04be59e64362d

WORKDIR /work

# Build libsqlitefunctions.so.
COPY ./build/extension-functions.c .
RUN apt-get -qq update \
 && apt-get -qq --no-install-recommends install \
      gcc \
      libsqlite3-dev \
 && gcc -fPIC -shared -O3 extension-functions.c -o libsqlitefunctions.so -lm \
 && apt-get purge -qq gcc libsqlite3-dev \
 && rm -rf /var/lib/apt/lists/*

# Install dependencies.
RUN apt-get -qq update \
 && apt-get -qq --no-install-recommends install \
      brotli \
      curl \
      gzip \
      varnish \
 && rm -rf /var/lib/apt/lists/*
# https://github.com/curiousleo/datasette-leaflet-geojson/archive/1e402abeb77192e0b8d51504b46055f1e1b4cf4d.tar.gz
RUN pip install \
      datasette-cluster-map \
      datasette-graphql \
      datasette-vega \
 && true
RUN curl -sSL https://github.com/DarthSim/hivemind/releases/download/v1.1.0/hivemind-v1.1.0-linux-amd64.gz \
    | gunzip - >/usr/bin/hivemind \
 && chmod +x /usr/bin/hivemind

# The database files are renamed to .sqlite3 on purpose.
#
# The only way I can find in the documentation to load settings from a file is to use "Configuration
# directory mode", i.e. passing a directory to `datasette`.
#
# https://docs.datasette.io/en/0.60.1/settings.html#config-dir
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
COPY ./build/Marktstammdatenregister.db.br ./Marktstammdatenregister.sqlite3.br
COPY ./build/Procfile .
COPY ./build/datasette.sh .
COPY ./build/default.vcl .
COPY ./build/metadata.yaml .
COPY ./build/settings.json .
COPY ./build/varnish.sh .

EXPOSE 8080

CMD hivemind
