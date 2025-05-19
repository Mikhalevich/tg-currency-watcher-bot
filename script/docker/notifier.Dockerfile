FROM golang:1.24-alpine3.21 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -a -installsuffix cgo -ldflags="-w -s" -o ./bin/notifier cmd/notifier/main.go

FROM alpine:3.21

EXPOSE 8080

WORKDIR /app/

COPY --from=builder /app/bin/notifier /app/notifier
COPY --from=builder /app/config/config-notifier.yaml /app/config-notifier.yaml

ENTRYPOINT ["./notifier", "-config", "config-notifier.yaml"]
