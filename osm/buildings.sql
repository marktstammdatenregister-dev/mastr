-- Import PostgreSQL COPY format
create virtual table buildings_import using VirtualText(
    'boundaries.pg',
    'UTF-8',
    0,
    POINT,
    DOUBLEQUOTE,
    TAB
);

-- Cast EWKB-encoded GEOMETRYCOLLECTION to Multipolygon
create table buildings (
    geometry MULTIPOLYGON not null,
    tags text --not null
    --name text generated always as (json_extract(tags, '$.name')) virtual not null,
    --admin_level integer generated always as (json_extract(tags, '$.admin_level')) virtual
);

insert into buildings
select
    CastToMultipolygon(GeomFromEWKB(COL001)),
    COL002
from
    buildings_import
where
    CastToMultipolygon(GeomFromEWKB(COL001)) is not null;

insert into buildings
select
    CastToMultipolygon(BuildArea(CastToMultilinestring(GeomFromEWKB(COL001)))),
    COL002
from
    buildings_import
where
    CastToMultipolygon(GeomFromEWKB(COL001)) is null;

-- Drop import table
drop table buildings_import;

-- Optimize
analyze;
vacuum;
