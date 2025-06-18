# Multi-stage Dockerfile for mkcd
# Build stage
FROM golang:1.24.4-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o mkcd .

# Final stage
FROM scratch

# Import from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Copy the binary
COPY --from=builder /build/mkcd /usr/local/bin/mkcd

# Create a non-root user
USER 65534:65534

# Set entrypoint
ENTRYPOINT ["/usr/local/bin/mkcd"]

# Default command
CMD ["--help"]

# Metadata
LABEL org.opencontainers.image.title="mkcd"
LABEL org.opencontainers.image.description="Enterprise directory creation and workspace initialization tool"
LABEL org.opencontainers.image.vendor="mochajutsu"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.source="https://github.com/mochajutsu/mkcd"
LABEL org.opencontainers.image.documentation="https://github.com/mochajutsu/mkcd/blob/main/README.md"
