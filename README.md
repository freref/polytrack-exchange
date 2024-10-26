# Polytrack Exchange

A central place for people to find tracks, leaderboard, stats and more for the game [PolyTrack](https://www.kodub.com/apps/polytrack). The website is written in Golang and HTMX and uses the standard librbary for templating.

## Requirements

- **Standalone Tailwind CLI** compiled with daisyUI, instructions to compile can be found [here](https://github.com/tailwindlabs/tailwindcss/discussions/12294#discussioncomment-8268378). A precompiled solution can be found [here](https://github.com/dobicinaitis/tailwind-cli-extra).
- **Docker**
- **Make**

## Instructions
### Build 
#### Create `.env` file

```
DB_HOST=db
DB_USER=user
DB_PASSWORD=password
DB_NAME=name
DB_PORT=port
DATABASE_URL=url
```

#### Build the project using make
```sh
  make build
```
