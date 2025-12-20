# Build stage
FROM golang:1.25-alpine AS builder

# Add metadata
LABEL org.opencontainers.image.title="Trigra - Kubernetes GitOps Controller"
LABEL org.opencontainers.image.description="Lightweight GitOps controller for Kubernetes clusters"
LABEL org.opencontainers.image.authors="Taiwrash"
LABEL org.opencontainers.image.source="https://github.com/Taiwrash/trigra"
LABEL org.opencontainers.image.licenses="MIT"

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s -X main.Version=1.0.0 -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -a -installsuffix cgo \
    -o trigra ./cmd/trigra

# Final stage - minimal runtime image
FROM alpine:latest

# Add metadata to final image
LABEL org.opencontainers.image.title="Trigra - Kubernetes GitOps Controller"
LABEL org.opencontainers.image.description="Lightweight GitOps controller for Kubernetes clusters"
LABEL org.opencontainers.image.authors="Taiwrash"
LABEL org.opencontainers.image.source="https://github.com/Taiwrash/trigra"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.version="1.0.0"

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata && \
    addgroup -g 1000 trigra && \
    adduser -D -u 1000 -G trigra trigra

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /build/trigra /bin/trigra

# Change ownership
RUN chown -R trigra:trigra /app

# Switch to non-root user
USER trigra

# Expose port
EXPOSE 8082

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8082/health || exit 1

# Run the application
ENTRYPOINT ["/bin/trigra"]
