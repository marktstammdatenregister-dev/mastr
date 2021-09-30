-- Import buildings.
create virtual table imported using VirtualText(
    'buildings.tsv',
    'UTF-8',
    0,
    POINT,
    DOUBLEQUOTE,
    TAB
);

select 'Number of buildings:           ' || count(*) from imported;

-- Copy into 'buildings' table and drop virtual table.
create table buildings (geometry MULTIPOLYGON not null);
insert into buildings select SetSRID(MultiPolygonFromText(COL001), 4326) from imported;
drop table imported;

-- Index on area.
create index idx_buildings_area on buildings (Area(Transform(geometry, 25832)));
select 'Number of buildings <1000 sqm: ' || count(*) from buildings where Area(Transform(geometry, 25832)) < 1000;
delete from buildings where Area(Transform(geometry, 25832)) < 1000;
select 'Minimum area [square meters]:  ' || min(Area(Transform(geometry, 25832))) from buildings;
select 'Maximum area [square meters]:  ' || max(Area(Transform(geometry, 25832))) from buildings;

-- Create spatial index.
select RecoverGeometryColumn(
    'buildings',
    'geometry',
    4326,
    'MULTIPOLYGON'
);
select CreateSpatialIndex('buildings', 'geometry');
select UpdateLayerStatistics('buildings', 'geometry');
