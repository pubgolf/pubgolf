# api

This directory contains the PubGolf API server, which provides a gRPC-web API for all clients (web, iOS, admin), as well as some platform-specific HTTP endpoints for the web app (including serving the web app assets via reverse proxy).

## Setup

Make sure you've installed the following:

* Go 1.17+
* Docker / Docker Compose
* Doppler

Configure `doppler` at the repo root to use the `api-server` project and `dev` environment by running `doppler setup`. 

If you are running the web app frontend against your local API server (instead of the remote staging environment), configure `doppler` in the `web-app` directory to use the `web-app` project and the `dev_local_api` environment (NOT `dev`).

## Development

### Running

Run the server using the following command:

```bash
doppler run -- bin/start 
```

With the server running, you should be able to see plaintext output at [localhost:5000/health-check](http://localhost:5000/health-check) via the browser.

If you want to run latest dev version of the web frontend as well, navigate to the `web-app` directory in another terminal tab and start the dev server on port 3000 (`doppler run -- npm run dev`).
