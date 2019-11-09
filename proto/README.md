# Proto Definitions

Definition files for protobuf and and GRPC. Used by the API, marketing site, web-app and iOS client.

## Setup / Pre-requisites

Install the following:

* Homebrew
* protoc (`brew install protobuf` or [download from here](https://github.com/protocolbuffers/protobuf/releases))
* [protoc-gen-grpc-web plugin](https://github.com/grpc/grpc-web/releases)

## Development

Instructions for compiling the protos for use in the other build targets.

### Web App

```
cd ../web-app
protoc -I=../proto ../proto/*.proto --js_out=import_style=commonjs:./src/proto --grpc-web_out=import_style=commonjs,mode=grpcwebtext:./src/proto
```
```
