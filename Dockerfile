FROM golang:1.23.5-alpine AS builder

WORKDIR /app

# Copy go module files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o dill-monitor ./cmd/server

# Use a smaller image for the final stage
FROM alpine:latest

WORKDIR /app

# Install CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Copy the binary from builder stage
COPY --from=builder /app/dill-monitor /app/dill-monitor

# Create config directory
RUN mkdir -p /app/config

# Create entrypoint script
RUN echo '#!/bin/sh' > /app/entrypoint.sh && \
    echo 'if [ ! -f /app/config/config.json ]; then' >> /app/entrypoint.sh && \
    echo '  echo "Creating default config.json"' >> /app/entrypoint.sh && \
    echo '  echo "{\"addresses\":[]}" > /app/config/config.json' >> /app/entrypoint.sh && \
    echo 'fi' >> /app/entrypoint.sh && \
    echo 'if [ ! -f /app/config/server_config.json ]; then' >> /app/entrypoint.sh && \
    echo '  echo "Creating default server_config.json"' >> /app/entrypoint.sh && \
    echo '  echo "{\"metricsPort\":9090,\"logLevel\":\"info\",\"host\":\"0.0.0.0\"}" > /app/config/server_config.json' >> /app/entrypoint.sh && \
    echo 'fi' >> /app/entrypoint.sh && \
    echo 'exec /app/dill-monitor -config=/app/config/config.json -server-config=/app/config/server_config.json "$@"' >> /app/entrypoint.sh && \
    chmod +x /app/entrypoint.sh

# Create volume for persistent data
VOLUME ["/app/config"]

# Set the entry point to our script
ENTRYPOINT ["/app/entrypoint.sh"]
