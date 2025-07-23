{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };
  outputs = {
    self,
    nixpkgs,
    flake-utils,
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = nixpkgs.legacyPackages.${system};
        buildPkgs = with pkgs; [
          pkg-config
          scdoc
        ];
        libPkgs = with pkgs; [
          libmpdclient
        ];
        devPkgs = with pkgs; [
          just
          go
          golangci-lint
          vhs
        ];
      in {
        packages.default = pkgs.buildGoModule {
          pname = "mpcube";
          version = "1.0.0";
          src = ./.;
          vendorHash = "sha256-VZuTMhjFEGWHhBJ2pukiIyQrHSo3LAB/2Ig9/5OsGjM=";
          nativeBuildInputs = buildPkgs;
          buildInputs = libPkgs;
          postInstall = ''
            mkdir -p $out/share/man/man1
            scdoc < mpcube.1.scd > $out/share/man/man1/mpcube.1
          '';
        };
        devShell = pkgs.mkShell {
          nativeBuildInputs = buildPkgs ++ [pkgs.go];
          buildInputs = libPkgs ++ devPkgs;
        };
      }
    );
}
