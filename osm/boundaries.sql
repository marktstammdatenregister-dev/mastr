--create table boundaries_import (geom_wkb_hex text, tags text);
--.mode tabs
--.import boundaries.pg boundaries_import
--alter table boundaries_import add column geom GEOMETRY;
--update boundaries_import set geom = GeomFromEWKB(geom_wkb_hex);

--select ImportGeoJSON('boundaries.geojson', 'boundaries_import');
create virtual table boundaries_import using VirtualGeoJSON('boundaries.geojson');

-- Cast GEOMETRYCOLLECTION to Multipolygon

create table boundaries (name text not null, geometry MULTIPOLYGON not null);
insert into boundaries select name, CastToMultipolygon(geometry) from boundaries_import where CastToMultipolygon(geometry) is not null;
insert into boundaries select name, CastToMultipolygon(BuildArea(CastToMultilinestring(geometry))) from boundaries_import where CastToMultipolygon(geometry) is null;
select name from boundaries where geometry is null;

--drop table boundaries_import;
select DiscardGeometryColumn('boundaries_import', 'geometry');
select DropVirtualGeometry('boundaries_import');
