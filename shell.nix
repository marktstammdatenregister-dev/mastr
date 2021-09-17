{ pkgs ? import <nixpkgs> {} }:
pkgs.mkShell {
  nativeBuildInputs = with pkgs; [
    axel
    brotli
    flyctl
    gdal_2
    pup
    pypy3Packages.xmlschema
    rlwrap
    spatialite_tools
    sqlite
    unzip
  ];
}
