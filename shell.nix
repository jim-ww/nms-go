{
  pkgs ? import <nixpkgs> { },
}:
with pkgs;
mkShell {
  buildInputs = [
    go
    tailwindcss
    sqlite
    air
  ];

  shellHook = ''
    alias test='go test ./...'
    alias test-verbose='go test -v -cover ./...'
    alias build-css='tailwindcss build -o web/static/tailwind.css'
    alias run='build-css && air'
  '';
}
