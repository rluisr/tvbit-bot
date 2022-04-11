# syntax = docker/dockerfile:1.3-labs

FROM golang:1-alpine as builder
WORKDIR /go/src/tvbit-bot
COPY . .
RUN apk --no-cache add git openssh build-base
RUN cd cmd && go build -o app .

FROM alpine as production
EXPOSE 8080
RUN <<EOF
    apk add --no-cache ca-certificates libc6-compat tzdata \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime
    rm -rf /var/cache/apk/*
EOF
ENV TZ="Asia/Tokyo"
COPY --from=builder /go/src/tvbit-bot/cmd/app /app
ENTRYPOINT ["/app"]

