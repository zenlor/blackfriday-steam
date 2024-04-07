{
  description = "markdown to steam markup library and command";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-23.11";

    flakelight.url = "github:nix-community/flakelight";
    flakelight.inputs.nixpkgs.follows = "nixpkgs";
  };
  outputs = { self, flakelight, nixpkgs, ...}@inputs: flakelight ./. {
    inherit inputs;
    systems = nixpkgs.lib.systems.flakeExposed;

    devShell.packages = pkgs: [
      pkgs.go
      pkgs.gopls
      pkgs.nixpkgs-fmt
      pkgs.nixd
    ];

    package = { pkgs, lib, buildGoModule, ... }: buildGoModule {
      name = "md2steam";
      src = ./.;
      nativeBuildInputs = [ pkgs.go ];
      vendorHash = "sha256-LI2aTn0MH4x1sRN6wiihHn+fvNJ3KbwbRRIEQQMFE3s=";
      meta = {
        platforms = lib.platforms.all;
      };
    };

    formatters = {
      "*.go" = "go fmt";
    };
  };
}
