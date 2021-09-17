-- Go fast!
PRAGMA synchronous=OFF;

-- Precompute area.
ALTER TABLE multipolygons ADD COLUMN area REAL;
UPDATE multipolygons SET area = Area(Transform(GEOMETRY, 25832));

-- Create area index: we want to be able to order by area.
CREATE INDEX idx_multipolygons_area ON multipolygons(area);

-- Create operator index: we want to be able to filter by operator.
CREATE INDEX idx_multipolygons_operator ON multipolygons(operator);

-- Create osm_id index: we want to be able to look up by osm_id.
CREATE UNIQUE INDEX idx_multipolygons_osm_id ON multipolygons(osm_id);

-- Cleanup.
VACUUM;
