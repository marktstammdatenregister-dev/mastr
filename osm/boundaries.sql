-- Import PostgreSQL COPY format
create table boundaries_import (
    ewkb_hex text not null,
    tags text not null
);
.mode tabs
.import 'boundaries.pg' 'boundaries_import'
delete from boundaries_import where json_extract(tags, '$.name') is null;

-- Cast EWKB-encoded GEOMETRYCOLLECTION to Multipolygon
create table boundaries (
    geometry MULTIPOLYGON not null,
    tags text not null,
    name text generated always as (json_extract(tags, '$.name')) virtual not null,
    admin_level integer generated always as (json_extract(tags, '$.admin_level')) virtual
);

insert into boundaries
select
    CastToMultipolygon(GeomFromEWKB(ewkb_hex)),
    tags
from
    boundaries_import
where
    CastToMultipolygon(GeomFromEWKB(ewkb_hex)) is not null;

insert into boundaries
select
    CastToMultipolygon(BuildArea(CastToMultilinestring(GeomFromEWKB(ewkb_hex)))),
    tags
from
    boundaries_import
where
    CastToMultipolygon(GeomFromEWKB(ewkb_hex)) is null;

-- Drop import table
drop table boundaries_import;
vacuum;
