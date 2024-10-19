# Polytrack Exchange

A central place for people to find maps, leaderboard, records and more for the game [PolyTrack](https://www.kodub.com/apps/polytrack). The website is written in Golang and HTMX.

# Requirements

- **Standalone Tailwind CLI** compiled with daisyUI, instructions to compile can be found [here](https://github.com/tailwindlabs/tailwindcss/discussions/12294#discussioncomment-8268378). A precompiled solution can be found [here](https://github.com/dobicinaitis/tailwind-cli-extra).
- **Go**

# Instructions

## Generate tailwind.css

```sh
tailwindcss-extra -i ./static/css/input.css -o ./static/css/output.css --watch
```

## Build

```sh
docker compose up --build
```
