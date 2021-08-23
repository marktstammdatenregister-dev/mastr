import fire
import geopandas as gpd
import plotly.express as px


def area(input_file, output_file):
    df = gpd.read_file(input_file)
    df = df.set_crs("EPSG:4326")

    area = df.to_crs("ESRI:102013")["geometry"].area
    df["area"] = area

    df = df.sort_values(by=["area"], ascending=False)

    fig = px.choropleth(
        df,
        geojson=df,
        featureidkey="id",
        locations="id",
        color="area",
        #mapbox_style="open-street-map",
        #center={"lat": 48.801067, "lon": 8.324541},
    )
    with open(output_file, "w") as f:
        fig.write_html(f)

    print(df)


if __name__ == "__main__":
    fire.Fire()
