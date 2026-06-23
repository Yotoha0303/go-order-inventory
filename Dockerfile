# syntax=docker/dockerfile:1

FROM golang:1.25.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -o go-order-inventory ./cmd

FROM alpine:3.22

WORKDIR /app

RUN addgroup -S app && adduser -S app -G app

COPY --from=builder /app/go-order-inventory ./go-order-inventory
COPY config.yml ./config.yml

USER app

EXPOSE 8082

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
	CMD wget -qO- http://127.0.0.1:8082/ping || exit 1

STOPSIGNAL SIGTERM

CMD ["./go-order-inventory"]