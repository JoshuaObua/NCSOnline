# ── Stage 1: Build ────────────────────────────────────────────────
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

# Copy module definition first for layer caching
COPY go.mod ./
# go.sum is generated automatically if absent (-mod=mod)
COPY go.sum* ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -mod=mod -ldflags="-w -s" -o ncsms-backend ./cmd/server

# ── Stage 2: Run ──────────────────────────────────────────────────
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app/ncsms-backend /ncsms-backend

EXPOSE 8080

ENTRYPOINT ["/ncsms-backend"]
