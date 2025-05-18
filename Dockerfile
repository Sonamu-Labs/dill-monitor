FROM golang:1.21-alpine AS builder

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

# Handle configuration files
RUN mkdir -p /tmp/config-setup
COPY config/ /tmp/config-setup/
RUN if [ -f /tmp/config-setup/config.json ]; then cp /tmp/config-setup/config.json /app/config/; fi && \
    if [ -f /tmp/config-setup/server_config.json ]; then cp /tmp/config-setup/server_config.json /app/config/; fi && \
    rm -rf /tmp/config-setup

# Create volume for persistent data
VOLUME ["/app/config"]

# Set the entry point
ENTRYPOINT ["/app/dill-monitor"]
CMD ["-config=/app/config/config.json", "-server-config=/app/config/server_config.json"]
