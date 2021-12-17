# web-app

This directory contains the frontend web app for [pubgolf.co](https://pubgolf.co).

## Setup

Make sure you have the following installed:

- NVM
- Doppler

Configure `doppler` in the `web-app` directory to use the `web-app` project and the `dev` environment.

**Note:** The following isn't yet accurate. Adding for future reference.

This will point to the staging API server. If you are running the web app frontend against a local API server (instead of the remote staging environment), configure `doppler` in the `web-app` directory to use the `web-app` project and the `dev_local_api` environment (NOT `dev`).

## Development

Remember to run `nvm use` whenever you navigate into the repo.

### Running

Run locally with hot reloading:

```bash
doppler run -- npm run dev
```

### Build

Verify that production builds will be generated correctly:

```bash
doppler run -- npm run build
```

### Lint / Test

```bash
npm run check
npm run format
npm run lint
```
