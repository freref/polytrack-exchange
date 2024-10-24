# Makefile

all: build
	
build: download_htmx tailwind docker_up

download_htmx:
	mkdir -p static/js
	wget -O static/js/htmx.min.js https://unpkg.com/htmx.org@2.0.3/dist/htmx.min.js

tailwind_watch:
	tailwindcss-extra -i ./static/css/input.css -o ./static/css/output.css --watch

tailwind:
	tailwindcss-extra -i ./static/css/input.css -o ./static/css/output.css 

docker_up:
	docker compose up --build
