FROM golang:1.20 as builder

WORKDIR /torrent_client

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o torrent-client ./main.go

FROM alpine:latest

WORKDIR /torrent_client
COPY --from=builder /torrent_client/torrent-client .

EXPOSE 9090

CMD ["./torrent-client"]
