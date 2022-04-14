tvbit-bot
============

[![lint](https://github.com/rluisr/tvbit-bot/actions/workflows/lint.yml/badge.svg?branch=master)](https://github.com/rluisr/tvbit-bot/actions/workflows/lint.yml)

tvbit-bot is TradingView webhook handler for Bybit.

tvbit = T(rading)V(iew) (By)bit

Introduction
-------------

1. Run application
2. Set an alert with webhook and a message as JSON like below:

```json
{
  "is_test_net": true,
  "api_key": "6pNgDhklvGadvr8D95",
  "api_secret_key": "TRgXYMOOVUBzQVb7m38PZ1ze59UrQC1KTW1N",
  "order": {
    "symbol": "BTCUSDT",
    "type": "Market",
    "price": 0,
    "side": "Buy",
    "qty": 0.014,
    "tp": 0,
    "sl": 0
  }
}
```

see [tv.go](pkg/domain/tv.go) or [Bybit API Documentation](https://bybit-exchange.github.io/docs/linear/#:~:text=Transaction%20timestamp-,order,-How%20to%20Subscribe).

tvbit-bot conforms to Bybit API.

Setup
-----

You can change listen port with `PORT` environment variable.

### Docker

```shell
$ docker run ghcr.io/rluisr/tvbit-bot:latest --name tvbit-bot -p 8080:8080 -d
```

### Binary

1. Download a binary from [release page.](https://github.com/rluisr/tvbit-bot/releases)
2. `$ ./app`

tvbit-bot.hcloud.ltd
--------------------

URL: `https://tvbit-bot.hcloud.ltd/tv`

I am offering this application for public use.
But I am sure that I may betray you. You may use it for production operation, or you may try it only for testing.

Powered by [HCloud Ltd](https://hcloud.ltd)

### Terms of service

I accept no responsibility whatsoever.

Limitation
----------

tvbit-bot does not support to close/cancel positions now.

Welcome your PR.

Twitter [@rarirureluis](https://twitter.com/rarirureluis)
