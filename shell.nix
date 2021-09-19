{ pkgs ? import <nixpkgs> {} }:
pkgs.mkShell {
  nativeBuildInputs = with pkgs; [
    axel
    brotli
    docker-compose
    entr
    flyctl
    gdal_2
    libxml2
    pup
    pypy3Packages.xmlschema
    rlwrap
    spatialite_tools
    sqlite
    unzip
  ];
}
