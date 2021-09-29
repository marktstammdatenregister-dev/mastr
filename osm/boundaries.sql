-- Import PostgreSQL COPY format
create table boundaries_import (
    ewkb_hex text not null,
    tags text not null
);
.mode tabs
.import 'boundaries.pg' 'boundaries_import'

-- Delete invalid entries
select 'Total: ' || count(*) from boundaries_import;

select 'Invalid JSON: ' || count(*) from boundaries_import where json_valid(tags) = 0;
delete from boundaries_import where json_valid(tags) = 0;

select 'Missing "name" field: ' || count(*) from boundaries_import where json_extract(tags, '$.name') is null;
delete from boundaries_import where json_extract(tags, '$.name') is null;

-- Cast EWKB-encoded GEOMETRYCOLLECTION to Multipolygon
create table boundaries (
    tags text not null,
    name text generated always as (json_extract(tags, '$.name')) virtual not null,
    admin_level integer generated always as (json_extract(tags, '$.admin_level')) virtual
);
select AddGeometryColumn('boundaries', 'geometry', 4326, 'MULTIPOLYGON', 'XY', 1);

insert into boundaries (geometry, tags)
select
    CastToMultipolygon(GeomFromEWKB(ewkb_hex)),
    tags
from
    boundaries_import
where
    CastToMultipolygon(GeomFromEWKB(ewkb_hex)) is not null;

insert into boundaries (geometry, tags)
select
    CastToMultipolygon(BuildArea(CastToMultilinestring(GeomFromEWKB(ewkb_hex)))),
    tags
from
    boundaries_import
where
    CastToMultipolygon(GeomFromEWKB(ewkb_hex)) is null;

-- Create index for queries of the form 'select geometry from boundaries where name = ?'
select CreateSpatialIndex('boundaries', 'geometry');
create index idx_boundaries_name on boundaries (name);

-- Drop import table
drop table boundaries_import;

-- Optimize
analyze;
vacuum;
