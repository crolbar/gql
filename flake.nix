{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };
  };

  outputs = inputs:
    inputs.flake-parts.lib.mkFlake {inherit inputs;} {
      systems = ["x86_64-linux" "aarch64-linux"];

      perSystem = {system, ...}: let
        pkgs = import inputs.nixpkgs {inherit system;};
      in {
        devShells.default = pkgs.mkShell {
          nativeBuildInputs = with pkgs; [
            go
            gopls
          ];
        };

        packages = {
          default = pkgs.buildGoModule {
            pname = "gql";
            version = "0.1";

            src = ./.;

            checkPhase = ''
              go test -v -short ./...
            '';

            vendorHash = "sha256-yNBzWcSZzb+yywLaOtBd9sS0DP8XfkUAsnpBykFp5jI=";
          };
        };
      };
    };
}
