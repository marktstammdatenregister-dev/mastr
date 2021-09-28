-- Import PostgreSQL COPY format
create table buildings_import (
    ewkb_hex text not null,
    tags text not null
);
.mode tabs
.import 'buildings.pg' 'buildings_import'

-- Delete invalid entries
select 'Total: ' || count(*) from buildings_import;

select 'Invalid JSON: ' || count(*) from buildings_import where json_valid(tags) = 0;
delete from buildings_import where json_valid(tags) = 0;

-- Cast EWKB-encoded GEOMETRYCOLLECTION to Multipolygon
create table buildings (
    geometry MULTIPOLYGON not null,
    tags text not null
    --name text generated always as (json_extract(tags, '$.name')) virtual not null,
    --admin_level integer generated always as (json_extract(tags, '$.admin_level')) virtual
);

insert into buildings
select
    CastToMultipolygon(GeomFromEWKB(ewkb_hex)),
    tags
from
    buildings_import
where
    CastToMultipolygon(GeomFromEWKB(ewkb_hex)) is not null;

insert into buildings
select
    CastToMultipolygon(BuildArea(CastToMultilinestring(GeomFromEWKB(ewkb_hex)))),
    tags
from
    buildings_import
where
    CastToMultipolygon(GeomFromEWKB(ewkb_hex)) is null;

-- Drop import table
drop table buildings_import;

-- Optimize
analyze;
vacuum;
