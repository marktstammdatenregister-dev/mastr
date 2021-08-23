import fire
import geopandas as gpd
import geoplot as gplt
import logging
import matplotlib.pyplot as plt


def preprocess(input_osm, output_feather):
    import pyrosm

    logging.info("Loading input file ...")
    df = pyrosm.OSM(input_osm).get_buildings()

    logging.info("Calculating area ...")
    df["area"] = df.to_crs("ESRI:102013")["geometry"].area

    logging.info("Sorting by area ...")
    df.sort_values(by=["area"], ascending=False, inplace=True)

    logging.info("Saving as feather ...")
    df.to_feather(output_feather, compression="lz4")

    logging.info("Done.")
    print(df)


def largest(input_feather, output_geojson):
    logging.info("Loading input file ...")
    df = gpd.read_feather(input_feather)
    df = df[:100]
    # df = df.set_crs("EPSG:4326")

    logging.info("Saving as GeoJSON ...")
    df.to_file(output_geojson, driver="GeoJSON")
    # print(df[["area", "addr:street", "addr:housenumber", "building"]])


def _area(input_file, output_file):
    df = gpd.read_file(input_file)
    df = df.set_crs("EPSG:4326")

    area = df.to_crs("ESRI:102013")["geometry"].area
    df["area"] = area
    df.sort_values(by=["area"], ascending=False, inplace=True)

    ax = gplt.webmap(df["geometry"], zoom=16, figsize=(24, 18))
    gplt.choropleth(df["geometry"], hue=df["area"], legend=True, cmap="Oranges", ax=ax)

    plt.savefig(output_file)

    print(df[["area", "addr:street", "addr:housenumber", "building"]])


if __name__ == "__main__":
    logging.basicConfig(level=logging.INFO, format="%(asctime)s %(message)s")
    fire.Fire()
