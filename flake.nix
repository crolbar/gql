{
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs = inputs: let
    systems = ["x86_64-linux" "aarch64-linux"];
    forEachSystem = inputs.nixpkgs.lib.genAttrs systems;
    pkgsFor = inputs.nixpkgs.legacyPackages;
  in {
    devShells = forEachSystem (system: let
      pkgs = pkgsFor.${system};
    in {
      default = pkgs.mkShell {
        nativeBuildInputs = with pkgs; [
          go
          gopls
        ];
      };
    });

    packages = forEachSystem (system: {
      default = pkgsFor.${system}.callPackage ./package.nix {};
    });
  };
}
