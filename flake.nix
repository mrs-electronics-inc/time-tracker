{
  description = "A simple CLI tool to track time spent on different projects and tasks";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        
        # Go module attributes
        goModule = pkgs.buildGoModule {
          pname = "time-tracker";
          # We don't have real versions yet
          version = "alpha";
          
          src = ./src;

          goPackagePath = "github.com/mrs-electronics-inc/time-tracker";
          vendorHash = "sha256-bCeLxvAURw4A0T/jG2rVLe7fXEYRL+s3NCQ5iZvdEOM=";
          
          meta = with pkgs.lib; {
            description = "A simple CLI tool to track time spent on different projects and tasks";
            homepage = "https://github.com/mrs-electronics-inc/time-tracker";
            license = licenses.mit;
            platforms = platforms.all;
          };
        };
      in
      {
        packages.default = goModule;
        packages.time-tracker = goModule;

        apps.default = {
          type = "app";
          program = "${goModule}/bin/time-tracker";
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
          ];

          shellHook = ''
            export GOPATH=$PWD/.go
            export PATH=$GOPATH/bin:$PATH
          '';
        };
      });
}
