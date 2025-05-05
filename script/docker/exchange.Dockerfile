FROM golang:1.24-alpine3.21 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -a -installsuffix cgo -ldflags="-w -s" -o ./bin/exchange cmd/exchange/main.go

FROM alpine:3.21

WORKDIR /app/

COPY --from=builder /app/bin/exchange /app/exchange
COPY --from=builder /app/config/config-exchange.yaml /app/config-exchange.yaml

ENTRYPOINT ["./exchange", "-config", "config-exchange.yaml"]
