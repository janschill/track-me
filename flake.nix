{
  description = "Flake for the track-me project";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }: flake-utils.lib.eachDefaultSystem (system:
    let
      pkgs = import nixpkgs { inherit system; };
    in
    {
      devShell = pkgs.mkShell {
        name = "track-me-dev-shell";

        buildInputs = [
          pkgs.go
          pkgs.sqlite
        ];

        shellHook = ''
          export GOBIN="$PWD/bin"
          export PATH="$GOBIN:$PATH"
          echo "Development environment for track-me is ready!"
        '';
      };
    }
  );
}
