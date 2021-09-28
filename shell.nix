{ pkgs ? (import (import ./nix/sources.nix).nixpkgs {})}:
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
    niv
    pup
    python3Packages.shapely
    rlwrap
    spatialite_tools
    sqlite
    unzip
    yj
  ];
}
