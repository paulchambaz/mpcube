{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-23.11";
    fenix = {
      url = "github:nix-community/fenix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { self, nixpkgs, fenix }:
    let
    pkgs = nixpkgs.legacyPackages.x86_64-linux;
    fenixPkgs = fenix.packages.x86_64-linux;

    manifest = (pkgs.lib.importTOML ./Cargo.toml).package;

    buildPkgs = with pkgs; [
      scdoc
      libmpdclient
    ];

    libPkgs = with pkgs; [
    ];

    devPkgs = with pkgs; [
      just
      cargo-tarpaulin
      vhs
    ];

    environment = with pkgs; {
    };

    makePkgConfigPath = libPkgs: pkgs.lib.concatStringsSep ":" (map (pkg: "${pkg.dev}/lib/pkgconfig") libPkgs);

    rust-toolchain = fenixPkgs.fromToolchainFile {
      file = ./rust-toolchain.toml;
      sha256 = "sha256-7QfkHty6hSrgNM0fspycYkRcB82eEqYa4CoAJ9qA3tU=";
    };

    rustPackage = {
      pname = manifest.name;
      version = manifest.version;
      src = self;

      cargoLock.lockFile = ./Cargo.lock;
      cargoLock.outputHashes = {
        "mpd-0.1.0" = "sha256-YVatWNIfSd98shfzgdD5rJ40indfge/2bT54DtJIQ1k=";
      };

      nativeBuildInputs = [ rust-toolchain ];
      buildInputs = buildPkgs ++ libPkgs;

      configurePhase = ''
        export PATH=${pkgs.lib.makeBinPath buildPkgs}:$PATH
        export PKG_CONFIG_PATH=${makePkgConfigPath libPkgs}:$PKG_CONFIG_PATH
        export HOME=$(mktemp -d)
      '';

      postInstall = ''
        mkdir -p $out/share/man/man1
        scdoc < mpcube.1.scd > $out/share/man/man1/mpcube.1
      '';

      meta = with pkgs.lib; {
        description = manifest.description;
        homepage = manifest.homepage;
        license = licenses.gpl3Plus;
        maintainers = with maintainers; [ paulchambaz ];
      };
    } // environment;

    shell = {
      buildInputs = [ rust-toolchain ] ++ buildPkgs ++ libPkgs ++ devPkgs;
    } // environment;
    in
  {
    packages.x86_64-linux.mpcube = pkgs.rustPlatform.buildRustPackage rustPackage;
    defaultPackage.x86_64-linux = self.packages.x86_64-linux.mpcube;
    devShell.x86_64-linux = pkgs.mkShell shell;
  };
}
