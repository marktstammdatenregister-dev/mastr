{ pkgs ? import <nixpkgs> {} }:
pkgs.mkShell {
  nativeBuildInputs = with pkgs; [
    axel
    brotli
    docker-compose
    flyctl
    gdal_2
    go
    libxml2
    mustache-go
    pup
    rlwrap
    spatialite_tools
    sqlite
    unzip
    yj
  ];
}
