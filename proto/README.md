# proto

This directory contains the proto definitions for the API server, as well as tools for editing them.

All commands in this README should be run from the **repo root**, NOT the `proto/` subdirectory.

## Setup

Make sure you've installed the following:
* Go 1.17+
* NVM
* [buf](https://docs.buf.build/installation)

## Development

### Lint/Test

```bash
buf lint
buf breaking --against "https://github.com/pubgolf/pubgolf.git#branch=develop"
```

## Deployment

Run the following to generate code based on the latest proto definition:

```bash
proto/bin/clean-and-generate-protos
```

Be sure to check the generated code into git, as well as bump the version constants in `proto/versions/current/*`.

If you have made a breaking change to the proto (at the serialization level), also bump the minimum compatible version in `proto/versions/mincompatible/*` to force clients to refresh.
