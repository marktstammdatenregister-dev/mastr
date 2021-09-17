SELECT InitSpatialMetaData();

CREATE TABLE 'EinheitenSolar_points' (
  'EinheitMastrNummer' TEXT PRIMARY KEY
);
SELECT AddGeometryColumn('EinheitenSolar_points', 'point', 4326, 'POINT', 'XY', 1);

INSERT INTO EinheitenSolar_points
SELECT
  EinheitMastrNummer,
  SetSRID(MakePoint(Laengengrad, Breitengrad), 4326)
FROM EinheitenSolar
WHERE
  MakePoint(Laengengrad, Breitengrad) is not null;

SELECT CreateSpatialIndex('EinheitenSolar_points', 'point');
