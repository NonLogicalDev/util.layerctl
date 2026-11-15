{
  description = "Compose Git repositories and changes as ordered layers";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-26.05-darwin";

  outputs = { self, nixpkgs, ... }:
    let
      systems = [ "aarch64-darwin" "aarch64-linux" "x86_64-darwin" "x86_64-linux" ];
      forAllSystems = f: nixpkgs.lib.genAttrs systems (system: f nixpkgs.legacyPackages.${system});
    in
    {
      packages = forAllSystems (pkgs: {
        default = pkgs.callPackage ./nix/package.nix { revision = self.rev or "unknown"; };
        layerctl = pkgs.callPackage ./nix/package.nix { revision = self.rev or "unknown"; };
      });
      checks = forAllSystems (pkgs: {
        package = pkgs.callPackage ./nix/package.nix { revision = self.rev or "unknown"; };
      });
    };
}
