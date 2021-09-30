-- Import buildings.
create virtual table imported using VirtualText(
    'buildings.tsv',
    'UTF-8',
    0,
    POINT,
    DOUBLEQUOTE,
    TAB
);

-- Copy into 'buildings' table and drop virtual table.
create table buildings (geometry MULTIPOLYGON not null);
insert into buildings select SetSRID(MultiPolygonFromText(COL001), 4326) from imported;
drop table imported;

-- Drop "small" buildings.
--delete from buildings where Area(Transform(geometry, 25832)) < 1000;
create index idx_buildings_area on buildings (Area(Transform(geometry, 25832)));

-- Create spatial index.
select RecoverGeometryColumn(
    'buildings',
    'geometry',
    4326,
    'MULTIPOLYGON'
);
select CreateSpatialIndex('buildings', 'geometry');
select UpdateLayerStatistics('buildings', 'geometry');

-- Optimize
analyze;
