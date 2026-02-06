# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o raindrop-mcp .

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/raindrop-mcp .

# Run as non-root user for security
RUN adduser -D -u 1000 mcp && chown -R mcp:mcp /app
USER mcp

ENTRYPOINT ["/app/raindrop-mcp"]
