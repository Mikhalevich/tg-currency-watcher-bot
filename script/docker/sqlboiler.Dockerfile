FROM golang:1.24-alpine3.21 AS builder

WORKDIR /app

RUN GOBIN=/app go install github.com/volatiletech/sqlboiler/v4@latest
RUN GOBIN=/app go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql@latest

FROM alpine:3.21

WORKDIR /app/

ENV PATH="$PATH:/app/"

COPY --from=builder /app/sqlboiler /app/sqlboiler
COPY --from=builder /app/sqlboiler-psql /app/sqlboiler-psql
COPY config/sqlboiler.toml /app/sqlboiler.toml

ENTRYPOINT ["./sqlboiler", "psql", "--config", "sqlboiler.toml"]
