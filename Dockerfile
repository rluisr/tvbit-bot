# syntax = docker/dockerfile:1.3-labs

ARG APP_NAME="tvbit-bot"

FROM golang:1 as builder

WORKDIR /go/src/${APP_NAME}

COPY go.* .
RUN --mount=type=cache,target=/root/.cache/go-build go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app cmd/main.go

FROM gcr.io/distroless/static-debian12 as production

LABEL maintainer="rluisr" \
  org.opencontainers.image.url="https://github.com/rluisr/tvbit-bot" \
  org.opencontainers.image.source="https://github.com/rluisr/tvbit-bot" \
  org.opencontainers.image.vendor="rluisr" \
  org.opencontainers.image.title="tvbit-bot" \
  org.opencontainers.image.description="TradingView webhook handler for Bybit." \
  org.opencontainers.image.licenses="AGPL"

COPY --from=builder /app /app
ENTRYPOINT ["/app"]
