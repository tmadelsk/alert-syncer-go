# syntax=docker/dockerfile:1

# --- Stage 1: Build the Go binary ---
FROM golang:1.20-alpine AS builder
WORKDIR /app

# Copy go.mod and go.sum first (if present) to leverage layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the binary for Linux
ENV CGO_ENABLED=0 GOOS=linux
RUN go build -ldflags="-s -w" -o alert-ingest-service .

# --- Stage 2: Runtime image ---
FROM alpine:latest
WORKDIR /app

# Install CA certificates (if your app needs external HTTP calls)
RUN apk add --no-cache ca-certificates

# Copy the built binary
COPY --from=builder /app/alert-ingest-service .

# Expose the HTTP API port (assuming default 8080)
EXPOSE 8080

# Use non-root user for security (optional)
RUN adduser -D appuser
USER appuser

ENTRYPOINT ["./alert-ingest-service"]