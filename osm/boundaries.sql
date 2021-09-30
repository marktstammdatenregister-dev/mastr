-- Import PostgreSQL COPY format
create virtual table imported using VirtualText(
    'boundaries.tsv',
    'UTF-8',
    0,
    POINT,
    DOUBLEQUOTE,
    TAB
);

select 'Number of boundaries: ' || count(*) from imported;
select 'Invalid JSON:         ' || count(*) from imported where json_valid(COL002) = 0;
select 'Missing "name" field: ' || count(*) from imported where json_valid(COL002) and json_extract(COL002, '$.name') is null;

-- Cast EWKB-encoded GEOMETRYCOLLECTION to Multipolygon
create table boundaries (
    geometry MULTIPOLYGON not null,
    tags text not null,
    name text generated always as (json_extract(tags, '$.name')) virtual not null,
    admin_level integer generated always as (json_extract(tags, '$.admin_level')) virtual
);
insert into boundaries
select
    SetSRID(MultiPolygonFromText(COL001), 4326) as geometry,
    COL002 as tags
from
    imported
where
    json_valid(tags)
    and json_extract(tags, '$.name') is not null;

-- Drop import table
drop table imported;

-- Set geometry
select RecoverGeometryColumn(
    'boundaries',
    'geometry',
    4326,
    'MULTIPOLYGON'
);

-- Index on name
create index idx_boundaries_name on boundaries (name);
