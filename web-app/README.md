# web-app

This directory contains the source code for the PubGolf promotional website, as well as the web-based admin UI.

Unlike the other sub-projects, commands in this README should be run within the `web-app/` subdirectory.

## Development

First, ensure the app server is running, (use `pubgolf-devctrl run api` in the repo root). Then, in another terminal in the `web-app/` subdirectory, run:

```sh
npm run dev
```

You can then access the app on [http://127.0.0.1:3000](http://127.0.0.1:3000). To change the port, update the dev configuration in Doppler.

## Lint/Check

```sh
npm run lint
npm run format
npm run check
```

## Test

```sh
npm run test:unit
npm run test:e2e
```

## Build

The following produces a production build and serves it on port `:4173`. Note that you will need to override `PUBGOLF_WEB_APP_UPSTREAM_HOST` when running the API to proxy this, as the normal `npm run dev` command exposes the web app on port `:5173`.

```sh
npm run build && npm run preview
```
