tvbit-bot
============

[![lint](https://github.com/rluisr/tvbit-bot/actions/workflows/lint.yml/badge.svg?branch=master)](https://github.com/rluisr/tvbit-bot/actions/workflows/lint.yml)

tvbit-bot is TradingView webhook handler for Bybit.

tvbit = T(rading)V(iew) (By)bit

Twitter [@rarirureluis](https://twitter.com/rarirureluis)

Introduction
-------------

1. Set an alert with webhook and a message as JSON like below:

```json
{
  "name": "alert name, description or something",
  "symbol": "BTCUSDT",
  "type": "Market",
  "price": "0",
  // If type is "Limit" set it as an int greater than 0
  "side": "Buy",
  "qty": "0.014",
  "tp": "0",
  // see below
  "sl": "{{high}}"
  // see below
}
```

more details, see [curl.txt](example/curl.txt)

### TP and SL

You need to set `tp` and `sl` as a string.

- `{{high}}` is an embedded value of TradingView, Also you can set any other TradingView's embedded values.
- Other methods, you can set as a percent like `"tp": "10%"` calculate from mark price.
- `"tp": "+40", "sl": "-20"` means, `TP: mark price + 40` and `SL: mark price - 20`.

see [tv.go](pkg/domain/tv.go)
or [Bybit API Documentation](https://bybit-exchange.github.io/docs/linear/#:~:text=Transaction%20timestamp-,order,-How%20to%20Subscribe)

Path
-----

| Path     | Method | Description             |
|----------|--------|-------------------------|
| /tv      | POST   | Create order            |


Setup
-----

You have to set environment variables

- [bybit](pkg/external/bybit/config.go)
- [mysql](pkg/external/mysql/config.go)

### Docker

`ghcr.io/rluisr/tvbit-bot:latest`

### Binary

[Release page.](https://github.com/rluisr/tvbit-bot/releases)

### MySQL

tvbit-bot saves the order history to MySQL.

Limitation
----------

tvbit-bot does not support to close/cancel positions, recommend to use TP/SL.

Welcome your PR.

Twitter [@rarirureluis](https://twitter.com/rarirureluis)
