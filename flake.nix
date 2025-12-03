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
        pname = "time-tracker";
        version = "alpha";
      in
      {
        packages.default = pkgs.buildGoModule {
          inherit pname version;
          src = self;
          
          # vendorHash locks Go module dependencies
          vendorHash = "sha256-2caz5wagKxYEBWkHpkdY3rv/K7Vvpqbt0DFK86N5oeY=";
          
          meta = with pkgs.lib; {
            description = "A simple CLI tool to track time spent on different projects and tasks";
            homepage = "https://github.com/mrs-electronics-inc/time-tracker";
            license = licenses.mit;
            platforms = platforms.all;
          };
        };

        packages.time-tracker = self.packages.${system}.default;

        apps.default = {
          type = "app";
          program = "${self.packages.${system}.default}/bin/time-tracker";
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            just
            docker
            docker-compose
          ];

          shellHook = ''
            export GOPATH=$PWD/.go
            export PATH=$GOPATH/bin:$PATH
          '';
        };
      });
}
