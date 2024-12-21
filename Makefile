tailwind-watch:
	tailwindcss build -o ./web/tailwind.css --minify --watch

templ-watch:
	templ generate -watch

dev:
	tailwindcss build -o ./web/tailwind.css --minify --watch & disown
	templ generate -watch & disown
	air

test:
	go test ./...
 
test-verbose:
	go test -v -cover ./...
