{
  pkgs ? import <nixpkgs> { },
}:
with pkgs;
mkShell {
  buildInputs = [
    go
    tailwindcss
    templ
    sqlite
    sqlc
    air
    go-migrate
  ];
}
