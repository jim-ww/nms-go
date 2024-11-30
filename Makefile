tailwind-watch:
	tailwindcss build -o ./static/tailwind.css --minify --watch

templ-watch:
	templ generate -watch -lazy

dev:
	tailwindcss build -o ./static/tailwind.css --minify --watch & disown
	templ generate -watch -lazy & disown
	air
