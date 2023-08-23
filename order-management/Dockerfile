FROM golang:1.20.6


WORKDIR /usr/src/app


COPY go.mod go.sum ./
RUN go mod dowload && go mod verify

COPY grpcapi/ ./grpcapi
COPY grpcclient/ ./grpcclient
COPY orderdb/ ./orderdb
COPY server/ ./server

RUN go build -v -o /usr/local/bin/app ./...


CMD ["app"]
