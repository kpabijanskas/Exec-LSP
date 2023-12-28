{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }: 
    flake-utils.lib.eachDefaultSystem (system:
      let
        name = "exec-lsp";
        version = "1.0.0";
        pkgs = nixpkgs.legacyPackages."${system}";
      in {
        packages.default = pkgs.buildGoModule {
          name = "${name}";

          src = ./.;

          vendorHash = "sha256-SUG7ldwOOWNKtrCrLgBFvwILTIy4kOZGfra7VfwJIXc=";

          meta = with pkgs.lib; {
            description = "Exec LSP is a very simple LSP server to execute commands.";
            homepage = "https://github.com/kpabijanskas/exec-lsp";
            license = licenses.gpl3;
          };
        };
      }
    );
}
