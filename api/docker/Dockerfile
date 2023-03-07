FROM golang:1.20

ADD ./api /go/src/github.com/pubgolf/pubgolf/api
ADD ./go.mod /go/src/github.com/pubgolf/pubgolf/go.mod
ADD ./go.sum /go/src/github.com/pubgolf/pubgolf/go.sum

WORKDIR /go/src/github.com/pubgolf/pubgolf/api
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/pubgolf-app-server github.com/pubgolf/pubgolf/api/cmd/pubgolf-app-server

FROM dopplerhq/cli:3.54
COPY --from=0 /go/bin/pubgolf-app-server .

CMD ["run", "--", "./pubgolf-app-server"]
