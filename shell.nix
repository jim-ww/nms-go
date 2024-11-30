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
  ];

  shellHook = ''
    alias tailwind-watch='tailwindcss build -o static/tailwind.css --watch'
    alias tailwind-build='tailwindcss build -o static/tailwind.css --minify'
    alias templ-watch='templ generate -watch -lazy'
    alias dev='tailwind-watch & disown && templ-watch & disown && air'
    alias test='go test ./...'
    alias test-verbose='go test -v -cover ./...'
    alias help='echo "$HELP"'
    export HELP='
      $ dev - starts server with air
      $ tailwind-build - rebuild css with tailwind
      $ tailwind-watch - rebuild css with tailwind and watch for changes
      $ test - run all tests
      $ test-verbose - run tests verbose + with coverage
      $ help - print this message'
    echo "$HELP"
  '';
}
