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
    alias build-css='tailwindcss build -o web/static/tailwind.css'
    alias run='build-css && air'
    alias test='go test ./...'
    alias test-verbose='go test -v -cover ./...'
    alias help='echo "$HELP"'
    export HELP='
      $ run - rebuilds c and starts server with air
      $ build-css - rebuild css with tailwind
      $ test - run all tests
      $ test-verbose - run tests verbose + with coverage
      $ help - print this message'
    echo "$HELP"
  '';
}
