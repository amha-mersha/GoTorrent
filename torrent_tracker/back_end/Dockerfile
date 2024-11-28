FROM golang:1.20 as builder

WORKDIR /back_end

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o tracker-backend ./main.go

FROM alpine:latest

WORKDIR /back_end
COPY --from=builder /back_end/tracker-backend .

COPY .env.example .env

EXPOSE 8080

CMD ["./tracker-backend"]
