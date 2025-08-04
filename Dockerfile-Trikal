FROM golang:1.24.5-alpine AS builder

RUN apk update && apk add --no-cache git upx
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o trikal cmd/main.go
RUN upx --best --lzma trikal

FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates
COPY .env .env
COPY --from=builder /app/trikal .
COPY --from=builder /app/rss.yaml ./

ENTRYPOINT ["./trikal"]



