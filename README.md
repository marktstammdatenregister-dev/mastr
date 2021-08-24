```
ogr2ogr -f SQLite out.sqlite karlsruhe-regbez-latest.osm.pbf -progress -dsco SPATIALITE=YES -gt 65536
```


Takes 18 seconds:
```
select
  Area(Transform(GEOMETRY, 25832)) as area,
  ogc_fid,
  osm_id,
  osm_way_id,
  name,
  type,
  aeroway,
  amenity,
  admin_level,
  barrier,
  boundary,
  building,
  craft,
  geological,
  historic,
  land_area,
  landuse,
  leisure,
  man_made,
  military,
  [natural],
  office,
  place,
  shop,
  sport,
  tourism,
  other_tags,
  GEOMETRY
from
  multipolygons
where
  "building" is not null
order by
  area desc
limit
  101
```

```
CREATE INDEX multipolygons_area ON multipolygons(Area(Transform(GEOMETRY, 25832)));
CREATE INDEX multipolygons_building ON multipolygons(building); // For building facet -- doesn't work
CREATE INDEX multipolygons_administrative ON multipolygons(name, boundary); // For subqueries
```

```
select Area(Transform(GEOMETRY, 25832)) as area, ogc_fid, name, building from multipolygons where "building" is not null order by area desc limit 101;
```

```
select
  Area(Transform(GEOMETRY, 25832)) as area,
  AsGeoJSON(GEOMETRY),
  ogc_fid,
  name,
  building,
  other_tags
from
  multipolygons
where
  "building" is not null
order by
  area desc
limit
  101
```

```
select
  name,
  AsGeoJSON(GEOMETRY)
from
  multipolygons
where
  boundary == "administrative"
limit 50
```

```
select
  AsGeoJSON(GEOMETRY)
from
  multipolygons
where
  name == "Gaggenau" AND boundary == "administrative"
limit 1
```

```
select
  AsGeoJSON(GEOMETRY),
  Area(Transform(GEOMETRY, 25832)) as area,
  ogc_fid,
  name,
  building,
  other_tags
from
  multipolygons,
  (
    select
      GEOMETRY as city_boundary
    from
      multipolygons
    where
      name == "Gernsbach"
      AND boundary == "administrative"
    limit
      1
  )
where
  building is not null
  and within(GEOMETRY, city_boundary)
order by
  area desc
limit
  20
```

```
select
  AsGeoJSON(GEOMETRY),
  Area(Transform(GEOMETRY, 25832)) as area,
  ogc_fid,
  name,
  building,
  other_tags
from
  multipolygons,
  (
    select
      GEOMETRY as city_boundary
    from
      multipolygons
    where
      name == "Gernsbach"
      AND boundary == "administrative"
    limit
      1
  )
where
  building is not null
  and building != "church"
  and within(GEOMETRY, city_boundary)
  and area > 1000
order by
  area asc
```

```
select
  count(ogc_fid),
  sum(area)
from
  (
    select
      ogc_fid,
      Area(Transform(GEOMETRY, 25832)) as area
    from
      multipolygons,
      (
        select
          GEOMETRY as city_boundary
        from
          multipolygons
        where
          name == "Durmersheim"
          AND boundary == "administrative"
        limit
          1
      )
    where
      building is not null
      and building != "church"
      and within(GEOMETRY, city_boundary)
      and area > 1000
  )
```

```
BEGIN TRANSACTION;
DROP TABLE lines;
DROP TABLE multilinestrings;
DROP TABLE other_relations;
DROP TABLE points;
DELETE FROM multipolygons WHERE building is null and boundary != "administrative";
END TRANSACTION;
VACUUM;
```

```
select
  AsGeoJSON(GEOMETRY),
  Area(Transform(GEOMETRY, 25832)) as area,
  ogc_fid,
  name,
  amenity,
  building,
  other_tags
from
  multipolygons,
  (
    select
      GEOMETRY as city_boundary
    from
      multipolygons
    where
      name == :city_name
      AND boundary == "administrative"
    limit
      1
  )
where
  building is not null
  and (amenity is null or amenity != "place_of_worship")
  and within(GEOMETRY, city_boundary)
  and area > cast(:minimum_area as int)
order by
  area desc
```

https://wiki.openstreetmap.org/wiki/Key:roof:shape

Aerial view:
https://github.com/digidem/leaflet-bing-layer/blob/gh-pages/index.html
https://gitlab.com/IvanSanchez/Leaflet.GridLayer.GoogleMutant/-/tree/master