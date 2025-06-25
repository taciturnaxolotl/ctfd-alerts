{
  description = "⛳ sending alerts for leaderboard changes and new challenges on any ctfd.io instance";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs }:
    let
      allSystems = [
        "x86_64-linux" # 64-bit Intel/AMD Linux
        "aarch64-linux" # 64-bit ARM Linux
        "x86_64-darwin" # 64-bit Intel macOS
        "aarch64-darwin" # 64-bit ARM macOS
      ];
      forAllSystems = f: nixpkgs.lib.genAttrs allSystems (system: f {
        pkgs = import nixpkgs { inherit system; };
      });
    in
    {
      packages = forAllSystems ({ pkgs }: {
        default = pkgs.buildGoModule {
          pname = "ctfd-alerts";
          version = "0.0.1";
          subPackages = [ "." ];  # Build from root directory
          src = self;
          vendorHash = null;
        };
      });

      devShells = forAllSystems ({ pkgs }: {
        default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            gotools
            go-tools
            (pkgs.writeShellScriptBin "ctfd-alerts-dev" ''
              go build -o ./bin/ctfd-alerts
              ./bin/ctfd-alerts "$@" || true
            '')
            (pkgs.writeShellScriptBin "ctfd-alerts-build" ''
              echo "Building ctfd-alerts binaries for all platforms..."
              mkdir -p $PWD/bin

              # Build for Linux (64-bit)
              echo "Building for Linux (x86_64)..."
              GOOS=linux GOARCH=amd64 go build -o $PWD/bin/ctfd-alerts-linux-amd64

              # Build for Linux ARM (64-bit)
              echo "Building for Linux (aarch64)..."
              GOOS=linux GOARCH=arm64 go build -o $PWD/bin/ctfd-alerts-linux-arm64

              # Build for macOS (64-bit Intel)
              echo "Building for macOS (x86_64)..."
              GOOS=darwin GOARCH=amd64 go build -o $PWD/bin/ctfd-alerts-darwin-amd64

              # Build for macOS ARM (64-bit)
              echo "Building for macOS (aarch64)..."
              GOOS=darwin GOARCH=arm64 go build -o $PWD/bin/ctfd-alerts-darwin-arm64

              # Build for Windows (64-bit)
              echo "Building for Windows (x86_64)..."
              GOOS=windows GOARCH=amd64 go build -o $PWD/bin/ctfd-alerts-windows-amd64.exe

              echo "All binaries built successfully in $PWD/bin/"
              ls -la $PWD/bin/
            '')
          ];

          shellHook = ''
            export PATH=$PATH:$PWD/bin
            mkdir -p $PWD/bin
          '';
        };
      });

      apps = forAllSystems ({ pkgs }: {
        default = {
          type = "app";
          program = "${self.packages.${pkgs.system}.default}/bin/ctfd-alerts";
        };
        ctfd-alerts-dev = {
          type = "app";
          program = toString (pkgs.writeShellScript "ctfd-alerts-dev" ''
            go build -o ./bin/ctfd-alerts ./main.go
            ./bin/ctfd-alerts $* || true
          '');
        };
        ctfd-alerts-build = {
          type = "app";
          program = "${self.devShells.${pkgs.system}.default.inputDerivation}/bin/ctfd-alerts-build";
        };
      });
    };
}
