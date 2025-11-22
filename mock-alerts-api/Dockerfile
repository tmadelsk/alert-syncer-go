# syntax=docker/dockerfile:1

# --- Stage 1: build the Go binary ---
FROM golang:1.20-alpine AS builder
WORKDIR /app

# Copy go.mod/go.sum first to leverage caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code
COPY . .

# Build the binary
# Assuming your `main.go` is in root of this module
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o mock-alerts-api .

# --- Stage 2: minimal runtime image ---
FROM alpine:latest
WORKDIR /app

# Copy the built binary from builder
COPY --from=builder /app/mock-alerts-api .

# Expose the port your service listens on; say 8081
EXPOSE 8081

# Use a non-root user for safety (optional)
RUN adduser -D appuser
USER appuser

# Entry point
ENTRYPOINT ["./mock-alerts-api"]