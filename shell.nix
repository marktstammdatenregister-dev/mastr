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
