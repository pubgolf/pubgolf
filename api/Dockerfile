FROM golang:1.14.1

RUN apt-get update && apt-get -y install unzip && apt-get clean

# Protobuf
ENV PB_VER 3.11.4
ENV PB_URL https://github.com/google/protobuf/releases/download/v${PB_VER}/protoc-${PB_VER}-linux-x86_64.zip
RUN mkdir -p /tmp/protoc && \
    curl -L ${PB_URL} > /tmp/protoc/protoc.zip && \
    cd /tmp/protoc && \
    unzip protoc.zip && \
    cp /tmp/protoc/bin/protoc /usr/local/bin && \
    cp -R /tmp/protoc/include/* /usr/local/include && \
    chmod go+rx /usr/local/bin/protoc && \
    cd /tmp && \
    rm -r /tmp/protoc

# protoc-gen-go
ENV GOPB_VER 1.3.5
RUN mkdir -p /tmp/protobuf && \
    curl -L https://github.com/golang/protobuf/archive/v${GOPB_VER}.zip > /tmp/protobuf/protobuf.zip && \
    cd /tmp/protobuf && \
    unzip protobuf.zip && \
    mkdir -p /go/src/github.com/golang/protobuf/ && \
    cp -r /tmp/protobuf/protobuf-${GOPB_VER}/* /go/src/github.com/golang/protobuf/ && \
    go install github.com/golang/protobuf/protoc-gen-go

ADD ./api /go/src/github.com/pubgolf/pubgolf/api
ADD ./proto /go/src/github.com/pubgolf/pubgolf/proto
ADD ./go.mod /go/src/github.com/pubgolf/pubgolf/go.mod
ADD ./go.sum /go/src/github.com/pubgolf/pubgolf/go.sum
WORKDIR /go/src/github.com/pubgolf/pubgolf/api
RUN bin/compile-protos
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/app github.com/pubgolf/pubgolf/api

FROM alpine:latest
COPY --from=0 /go/bin/app .

CMD ["./app"]
