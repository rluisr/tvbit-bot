# syntax = docker/dockerfile:1.3-labs

FROM golang:1-alpine as builder
ARG VERSION=0.0.0
WORKDIR /go/src/tvbit-bot
COPY . .
RUN apk --no-cache add git openssh build-base
RUN cd cmd && go build -ldflags "-X github.com/rluisr/tvbit-bot/pkg/external.version=${VERSION}" -o app .

FROM alpine as production
LABEL maintainer="rluisr" \
  org.opencontainers.image.url="https://github.com/rluisr/tvbit-bot" \
  org.opencontainers.image.source="https://github.com/rluisr/tvbit-bot" \
  org.opencontainers.image.vendor="rluisr" \
  org.opencontainers.image.title="tvbit-bot" \
  org.opencontainers.image.description="TradingView webhook handler for Bybit." \
  org.opencontainers.image.licenses="AGPL"
RUN <<EOF
    apk add --no-cache ca-certificates libc6-compat tzdata \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
    rm -rf /var/cache/apk/*
EOF
ENV TZ="Asia/Tokyo"
COPY --from=builder /go/src/tvbit-bot/cmd/app /app
ENTRYPOINT ["/app"]

