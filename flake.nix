{
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.11";
  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, nixpkgs, flake-utils, }: (flake-utils.lib.eachDefaultSystem (
    system: let
      pkgs = nixpkgs.legacyPackages.${system};
    in {
      packages.default = pkgs.buildGoModule {
        name = "plugin-buildah";
        src = ./.;
        vendorHash = "sha256-f5WANbL+aTgurdEpEFB3Ysh/TEQMjme/K/ckRgfZJWs=";
        version = "0.0.1";
      };
      devShells.default = pkgs.mkShell {
        packages = with pkgs; [
          go
          woodpecker-cli
        ];
      };
    }
  ));
}