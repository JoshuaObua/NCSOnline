# ── Stage 1: Build ────────────────────────────────────────────────
FROM golang:1.23-alpine AS builder

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

# ── Stage 2: Run ──────────────────────────────────────────────────
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app/ncsms-backend /ncsms-backend
COPY --from=builder /app/nsmis-worker /nsmis-worker

EXPOSE 8080

ENTRYPOINT ["/ncsms-backend"]
