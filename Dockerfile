
FROM golang:1.24.5-alpine AS builder
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -ldflags="-s -w" -o sarjan ./cmd/sarjan/main.go

FROM debian:bookworm-slim

WORKDIR /
COPY --from=builder /app/sarjan .
COPY .env .env
ENTRYPOINT ["/sarjan"]
