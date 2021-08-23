FROM neerteam/geopandas@sha256:a5f1de13268d3a71ea9f78fb88fd4499077bec8351b09e59403d1a6fe2bded4e

RUN pip3 install \
      fire \
      geoplot
