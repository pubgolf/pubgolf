# proto

This directory contains the proto definitions for the API server, as well as tools for editing them.

All commands in this README should be run from the **repo root**, NOT the `proto/` subdirectory.

## Lint/Test

```bash
buf lint
buf format -w
buf breaking --against "https://github.com/pubgolf/pubgolf.git#branch=develop"
```

## Deployment

Run the following to generate code based on the latest proto definition:

```bash
pubgolf-devctrl generate proto
```

Be sure to check the generated code into git for use with the web-app and API server. Generated client libs for mobile apps will be generated and published by the CI/CD process.
