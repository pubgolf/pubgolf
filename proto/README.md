# Proto Definitions

Definition files for protobuf and and GRPC. Used by the API, marketing site, web-app and iOS client.

## Setup / Pre-requisites

Install the following:

* Homebrew
* protoc (`brew install protobuf` or [download from here](https://github.com/protocolbuffers/protobuf/releases))
* web-app: [protoc-gen-grpc-web plugin](https://github.com/grpc/grpc-web/releases)
* api: protoc-gen-go (`go get -u protoc-gen-go`)

## Development

Instructions for compiling the protos for use in the other build targets.

### Web App

```
cd ../web-app
protoc -I=../proto ../proto/*.proto --js_out=import_style=commonjs:./src/proto --grpc-web_out=import_style=commonjs,mode=grpcwebtext:./src/proto
```

### API

```
cd ../api
protoc -I=../proto ../proto/*.proto --go_out=plugins=grpc:proto/pubgolf
```
