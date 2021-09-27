-- Create VirtualGeoJSON table
create virtual table boundaries_import using VirtualGeoJSON('boundaries.geojson');

-- Cast GEOMETRYCOLLECTION to Multipolygon
create table boundaries (
    name text not null,
    geometry MULTIPOLYGON not null
);

insert into boundaries
select
    name,
    CastToMultipolygon(geometry)
from
    boundaries_import
where
    CastToMultipolygon(geometry) is not null;

insert into boundaries
select
    name,
    CastToMultipolygon(BuildArea(CastToMultilinestring(geometry)))
from
    boundaries_import
where
    CastToMultipolygon(geometry) is null;

-- Drop VirtualGeoJSON table
select DropVirtualGeometry('boundaries_import');
