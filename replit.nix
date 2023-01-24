{ pkgs }: {
    deps = [
        pkgs.gosec
        pkgs.gotools
        pkgs.go
        pkgs.gopls
    ];
}