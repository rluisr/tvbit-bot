# syntax = docker/dockerfile:1.3-labs

ARG APP_NAME="tvbit-bot"
ARG VERSION=0.0.0

FROM golang:1 as builder

WORKDIR /go/src/${APP_NAME}
RUN go env -w GOMODCACHE=/root/.cache/go-build

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/root/.cache/go-build go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=0 GOOS=linux go build -ldflags "-X github.com/rluisr/tvbit-bot/pkg/external.version=${VERSION}" -o /app cmd/main.go

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
