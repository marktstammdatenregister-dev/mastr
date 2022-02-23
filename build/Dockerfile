#
# Run datasette.
#
FROM docker.io/datasetteproject/datasette:0.60.1@sha256:c2f93ba54eb49e299dd98180de0a8ca3cc1da0df05425a8e32f04be59e64362d
WORKDIR /work

# Build libsqlitemath.so.
RUN apt-get -qq update \
 && apt-get -qq --no-install-recommends install \
      gcc \
      libsqlite3-dev \
      wget \
 && wget -O sqlitemath.c https://sqlite.org/contrib/download/extension-functions.c?get=25 \
 && gcc -fPIC -shared sqlitemath.c -o libsqlitemath.so -lm \
 && apt-get purge -qq gcc libsqlite3-dev wget \
 && rm -rf /var/lib/apt/lists/*

# Install dependencies.
RUN apt-get -qq update \
 && apt-get -qq --no-install-recommends install \
      brotli \
      dumb-init \
 && rm -rf /var/lib/apt/lists/*
# https://github.com/curiousleo/datasette-leaflet-geojson/archive/1e402abeb77192e0b8d51504b46055f1e1b4cf4d.tar.gz
RUN pip install \
      datasette-cluster-map \
      datasette-graphql \
      datasette-vega \
 && true

RUN groupadd -r datasette && useradd --no-log-init -r -g datasette datasette
RUN chown datasette:datasette .
USER datasette:datasette

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
COPY ./build/metadata.yaml .
COPY ./build/settings.json .

ENTRYPOINT ["dumb-init", "--"]
EXPOSE 8080

# Use "." as the configuration directory. This loads *.db, metadata.yaml, settings.json, etc.
CMD ["sh", "-c", "brotli --rm --decompress --no-copy-stat *.br && datasette --port=8080 --host=0.0.0.0 --load-extension ./libsqlitemath.so --cors --immutable Marktstammdatenregister.sqlite3 ."]
