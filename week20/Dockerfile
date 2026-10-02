# ── Stage 1: Build ────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Cache dependencies first (layer caching)
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build a statically linked binary — no external deps needed
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server

# ── Stage 2: Runtime ──────────────────────────────────────────
FROM scratch

# Copy CA certs (for HTTPS calls from the container)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the binary only
COPY --from=builder /app/server /server

# Non-root user (numeric UID — scratch has no /etc/passwd)
USER 1001:1001

EXPOSE 8080

ENTRYPOINT ["/server"]