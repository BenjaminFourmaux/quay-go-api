# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS builder

WORKDIR /app

# Cache dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build the binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /out/quay-go-api ./main.go

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /out/quay-go-api /app/quay-go-api

EXPOSE 8080

ENTRYPOINT ["/app/quay-go-api"]