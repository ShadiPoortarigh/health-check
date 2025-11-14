FROM golang:1.25.3 AS builder

WORKDIR /app

COPY ../go.mod ../go.sum ./

COPY ../vendor ./vendor

ENV GOFLAGS="-mod=vendor"

COPY ../ ./

RUN go test ./...

RUN CGO_ENABLED=0 GOOS=linux go build -o healthcheck-app ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/healthcheck-app .

COPY ../sample-config.json .

EXPOSE 8080

CMD ["./healthcheck-app", "--config", "sample-config.json"]

