-- Go fast!
PRAGMA synchronous=OFF;

-- We're only interested in the 'multipolygons' table, so we drop all others.
DROP TABLE lines;
DROP TABLE multilinestrings;
DROP TABLE other_relations;
DROP TABLE points;

-- Precompute area.
CREATE INDEX multipolygons_area ON multipolygons(Area(Transform(GEOMETRY, 25832)));

-- Speed up queries of the form:
--
--     select
--       GEOMETRY as city_boundary
--     from
--       multipolygons
--     where
--       name == "Gernsbach"
--       AND boundary == "administrative"
--
CREATE INDEX multipolygons_administrative ON multipolygons(name, boundary);

-- Clean up and rebuild the database.
VACUUM;
