# ── Stage 1: Builder ──────────────────────────────────
FROM golang:1.25-alpine AS builder

# Install git (needed by some Go modules)
RUN apk add --no-cache git

WORKDIR /app

# Copy dependency files first — layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build a fully static binary — no external dependencies
RUN CGO_ENABLED=0 GOOS=linux go build -o taskboard .

# ── Stage 2: Final Image ───────────────────────────────
FROM alpine:3.19

WORKDIR /app

# Copy ONLY the binary from the builder stage
COPY --from=builder /app/taskboard .

EXPOSE 8080

CMD ["./taskboard"]
