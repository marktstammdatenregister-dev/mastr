import fire
import geopandas as gpd
import geoplot as gplt
import matplotlib.pyplot as plt


def area(input_file, output_file):
    df = gpd.read_file(input_file)
    df = df.set_crs("EPSG:4326")

    area = df.to_crs("ESRI:102013")["geometry"].area
    df["area"] = area

    df = df.sort_values(by=["area"], ascending=False)

    ax = gplt.webmap(df["geometry"], zoom=16, figsize=(24, 18))
    gplt.choropleth(df["geometry"], hue=df["area"], legend=True, cmap="Oranges", ax=ax)

    plt.savefig(output_file)

    print(df[["area", "addr:street", "addr:housenumber", "building"]])


if __name__ == "__main__":
    fire.Fire()
