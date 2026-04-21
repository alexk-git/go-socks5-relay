# Multi-stage build for SOCKS5 proxy
# Builder stage
FROM docker.io/library/golang:1.22-alpine AS builder

WORKDIR /app

# Copy go module files (if they exist)
COPY go.mod go.sum ./

# Initialize go module if go.mod doesn't exist
RUN if [ ! -f go.mod ]; then \
    echo "Initializing Go module..." && \
    go mod init go-socks5-relay && \
    go mod tidy; \
fi

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o socks5-proxy ./cmd/socks5-proxy

# Copy generation scripts
RUN mkdir -p /app/scripts
COPY scripts/generate-env.sh /app/scripts/
COPY scripts/entrypoint.sh /app/scripts/
RUN chmod +x /app/scripts/generate-env.sh /app/scripts/entrypoint.sh

# Final stage
FROM docker.io/library/alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/socks5-proxy .
# Copy scripts
COPY --from=builder /app/scripts/ ./scripts/
# Copy configuration example
COPY config/.env.example .

# Create volume for persistent configuration (optional)
# VOLUME /app/config



# Set entrypoint to generate .env and run proxy
ENTRYPOINT ["/app/scripts/entrypoint.sh"]

# Default command (can be overridden)
CMD []
