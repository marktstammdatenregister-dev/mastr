-- Go fast!
PRAGMA synchronous=OFF;

-- Precompute area.
ALTER TABLE multipolygons ADD COLUMN area INTEGER;
UPDATE multipolygons SET area = CAST(Area(Transform(GEOMETRY, 25832)) AS INTEGER);
CREATE INDEX idx_multipolygons_area ON multipolygons(area);
VACUUM;
