-- Import PostgreSQL COPY format
create virtual table imported using VirtualText(
    'buildings.txt',
    'UTF-8',
    0,
    POINT,
    DOUBLEQUOTE,
    TAB
);

-- Convert and copy into 'buildings'
create table buildings (geometry MULTIPOLYGON not null);
insert into buildings
    select SetSRID(MultiPolygonFromText(COL001), 4326)
    from imported;

-- Drop import table
drop table imported;

-- Create spatial index
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
vacuum;
