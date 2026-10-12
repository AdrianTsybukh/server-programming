{
  description = "Gio App Launcher Dev Environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};

        # 1. Gio compilation dependencies
        # Note: X11 libraries are under the 'xorg' namespace in NixOS
        buildDeps = with pkgs; [
          wayland
          libxkbcommon
          vulkan-headers
          xorg.libX11
          xorg.libXcursor
          xorg.libXfixes
        ];

        # 2. Runtime dependencies (needed for 'go run' to find libraries)
        runtimeDeps = with pkgs; [
          libGL
          vulkan-loader
        ];
      in
      {
        devShells.default = pkgs.mkShell {
          # We merge the lists using '++' so the shell has access to everything
          buildInputs = buildDeps ++ runtimeDeps ++ [ pkgs.go pkgs.pkg-config pkgs.gopls pkgs.go-tools];

          # 3. This environment variable tells the linker where to find .so files
          # The error "cannot coerce function" happened because you likely
          # removed 'runtimeDeps' from inside the parenthesis below.
          LD_LIBRARY_PATH = "${pkgs.lib.makeLibraryPath runtimeDeps}";
        };
      }
    );
}