# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install ca-certificates for HTTPS and git for private deps (optional)
RUN apk add --no-cache ca-certificates git

# Copy dependency files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /api ./cmd/api

# Runtime stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Non-root user
RUN adduser -D -g "" appuser

# Copy binary from builder
COPY --from=builder /api .

# Optional: persist SQLite DB when using DB_DRIVER=sqlite
# Use: docker run -v /app/data:/app/data ...
RUN mkdir -p /app/data && chown -R appuser:appuser /app
USER appuser

EXPOSE 3000

# Default: use env at runtime (e.g. docker run -e DB_DRIVER=postgres -e DATABASE_URL=...)
# For SQLite: mount a volume for /app/data and set DB_PATH=/app/data/data.db
ENTRYPOINT ["./api"]
