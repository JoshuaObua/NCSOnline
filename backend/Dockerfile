# ── Stage 1: Build ────────────────────────────────────────────────
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod .

# Skip module proxy and checksum DB — fetch directly from VCS.
# This avoids proxy.golang.org TLS issues on some networks.
ENV GOPROXY="direct"
ENV GONOSUMDB="*"
ENV GOFLAGS="-mod=mod"

# Download all deps (including transitives) before copying source.
RUN go mod download all

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOGC=50 \
    go build -p 2 -ldflags="-w -s" -o ncsms-backend ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOGC=50 \
    go build -p 2 -ldflags="-w -s" -o nsmis-worker ./cmd/worker
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GOGC=50 \
    go build -p 2 -ldflags="-w -s" -o ncs-ops-agent ./cmd/ops-agent

# ── Stage 2: Run ──────────────────────────────────────────────────
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata postgresql16-client rclone

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app/ncsms-backend /ncsms-backend
COPY --from=builder /app/nsmis-worker /nsmis-worker
COPY --from=builder /app/ncs-ops-agent /ncs-ops-agent

EXPOSE 8080

ENTRYPOINT ["/ncsms-backend"]
