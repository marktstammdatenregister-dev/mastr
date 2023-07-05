{ pkgs ? (import (import ./nix/sources.nix).nixpkgs {})}:
pkgs.mkShell {
  nativeBuildInputs = with pkgs; [
    axel
    brotli
    flyctl
    go_1_20
    libxml2
    mustache-go
    niv
    pigz
    pv
    rclone
    rlwrap
    shellcheck
    shfmt
    sqldiff
    sqlite
    sqlite-analyzer
    unzip
  ];
}
