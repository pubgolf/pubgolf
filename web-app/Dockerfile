FROM node:12.16.1

RUN apt-get -qq update && apt-get -qq install -y \
  unzip

WORKDIR /tmp
ENV PB_VER 3.8.0
RUN curl -sSL https://github.com/protocolbuffers/protobuf/releases/download/v${PB_VER}/\
protoc-${PB_VER}-linux-x86_64.zip -o protoc.zip && \
  unzip -qq protoc.zip && \
  cp ./bin/protoc /usr/local/bin/protoc

ENV GRPC_VER 1.0.7
RUN curl -sSL https://github.com/grpc/grpc-web/releases/download/${GRPC_VER}/\
protoc-gen-grpc-web-${GRPC_VER}-linux-x86_64 -o /usr/local/bin/protoc-gen-grpc-web && \
  chmod +x /usr/local/bin/protoc-gen-grpc-web

COPY ./proto /home/node/proto
COPY ./web-app /home/node/web-app

WORKDIR /home/node/web-app
RUN npm install -g npx && npm install
RUN npm run build

CMD ["node", "__sapper__/build/index.js"]
