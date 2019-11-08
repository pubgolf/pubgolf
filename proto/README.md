# Proto Definitions

Definition files for protobuf and and GRPC. Used by the API, marketing site, web-app and iOS client.

## Development

Instructions for compiling the protos for use in the other build targets.

### Web App

```
cd ../web-app
protoc -I=../proto pubgolf.proto --js_out=import_style=commonjs:./src/proto --grpc-web_out=import_style=commonjs,mode=grpcwebtext:./src/proto
```
