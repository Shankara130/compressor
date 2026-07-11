# =========================
# Builder Stage
# =========================
FROM golang:1.24.3-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -ldflags="-s -w" -o server ./cmd/server

RUN mkdir -p /app/tmp /app/output


# =========================
# Production Stage
# =========================
FROM alpine:3.20

# Runtime dependencies required by the optimizers:
#   ffmpeg      -> video optimization
#   ghostscript -> PDF optimization (provides the `gs` binary)
RUN apk add --no-cache ffmpeg ghostscript

WORKDIR /app

# Writable directories for uploads and processed output, owned by the
# non-root UID the server runs as (so uploads don't fail with EACCES).
RUN mkdir -p tmp output && chown -R 65532:65532 tmp output

COPY --from=builder /app/server .
COPY --from=builder /app/web ./web

USER 65532:65532

EXPOSE 8080

ENTRYPOINT ["/app/server"]
