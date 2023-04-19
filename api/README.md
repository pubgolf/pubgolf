# api

This directory contains the PubGolf API server, which provides a gRPC-web API for all clients (web, iOS, admin), as well as some platform-specific HTTP endpoints for the web app (including serving the web app assets via reverse proxy).

All commands in this README should be run from the **repo root**, NOT the `api/` subdirectory.

## Development

Run all dev commands (start server, run tests, generate mocks, etc) from the project root using the `pubgolf-devctrl` CLI tool.
