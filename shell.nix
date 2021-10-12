{ pkgs ? (import (import ./nix/sources.nix).nixpkgs {})}:
pkgs.mkShell {
  nativeBuildInputs = with pkgs; [
    axel
    brotli
    flyctl
    go
    libxml2
    mustache-go
    niv
    osmium-tool
    pup
    pv
    rlwrap
    spatialite_tools
    sqlite
    unzip
  ];
}
