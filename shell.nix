{ pkgs ? import <nixpkgs> {} }:
pkgs.mkShell {
  nativeBuildInputs = with pkgs; [
    axel
    brotli
    flyctl
    gdal_2
    rlwrap
    spatialite_tools
    sqlite
    unzip
  ];
}
