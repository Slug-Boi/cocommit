{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    prog_src = {
      flake = false;
      url = "github:Slug-Boi/cocommit";
    };
  };
  outputs = { self, nixpkgs, flake-utils, ... }@inputs: let
    forAllSys = flake-utils.lib.eachSystem flake-utils.lib.allSystems;

    APPNAME = "cocommit";
    appOverlay = final: prev: {
      # any pkgs overrides made here will be
      # inherited in the arguments of default.nix
      # because we used final.callPackage instead of prev.callPackage
      # i.e.
      # nodejs = prev.nodejs.overrideAttrs { name = "stinky"; };
      # would make it so that final.callPackage gives the altered nodejs
      ${APPNAME} = final.callPackage ./. { inherit (inputs) prog_src; };
    };
  in {
    overlays.default = appOverlay;
  } // (
    forAllSys (system: let
      pkgs = import nixpkgs { inherit system; overlays = [ appOverlay ]; };
    in{
      packages = {
        default = pkgs.${APPNAME};
      };
    })
  );
}
