{ pkgs ? (import (import ./nix/sources.nix).nixpkgs {})}:
pkgs.mkShell {
  nativeBuildInputs = with pkgs; [
    axel
    brotli
    docker-compose
    flyctl
    go
    libxml2
    mustache-go
    niv
    osmium-tool
    postgresql_13
    pup
    pv
    rlwrap
    spatialite_tools
    sqlite
    unzip
  ];
}
