FROM node:12.13.0

ENV NODE_ENV=production
EXPOSE ${WEB_APP_PORT}

RUN apt-get update && apt-get -y install unzip && apt-get clean

RUN npm install -g npx

# install protobuf
ENV PB_VER 3.6.0
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

# Get the source from GitHub
ENV GRPC_VER 1.0.7
RUN mkdir -p /tmp/grpc-web && \
    curl -L https://github.com/grpc/grpc-web/archive/protoc-gen-grpc-web-${GRPC_VER}-linux-x86_64 \
    > /tmp/grpc-web/protoc-gen-grpc-web && \
    cp -r /tmp/grpc-web/protoc-gen-grpc-web/* /usr/local/bin/protoc-gen-grpc-web && \
    chmod +x /usr/local/bin/protoc-gen-grpc-web && \
    rm -r /tmp/grpc-web

ADD ./web-app /go/src/github.com/escavelo/pubgolf/web-app
ADD ./proto /go/src/github.com/escavelo/pubgolf/proto
WORKDIR /go/src/github.com/escavelo/pubgolf/web-app

RUN npm install && npm run build

CMD ["node", "__sapper__/build"]
