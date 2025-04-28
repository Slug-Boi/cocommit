{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }@inputs: let
    forAllSys = flake-utils.lib.eachSystem flake-utils.lib.allSystems;

    APPNAME = "cocommit";
    appOverlay = final: prev: {
      ${APPNAME} = final.callPackage ./default.nix { 
        # Pass inputs as an argument to your package
        inherit (prev) lib fetchFromGitHub buildGoModule;
        # Or if you need all inputs:
        # inherit inputs;
      };
    };
  in {
    overlays.default = appOverlay;
  } // (
    forAllSys (system: let
      pkgs = import nixpkgs { 
        inherit system; 
        overlays = [ appOverlay ]; 
      };
    in {
      packages = {
        default = pkgs.${APPNAME};
      };
    })
  );
}
