tvbit-bot
============

[![lint](https://github.com/rluisr/tvbit-bot/actions/workflows/lint.yml/badge.svg?branch=master)](https://github.com/rluisr/tvbit-bot/actions/workflows/lint.yml)

[Japanese README](README_ja.md)

tvbit-bot is TradingView webhook handler for Bybit.

tvbit = T(rading)V(iew) (By)bit

Twitter [@rarirureluis](https://twitter.com/rarirureluis)

Introduction
-------------

1. Enable Webhook `https://<domain>/tv`
2. Set an alert with webhook and a message as JSON like below:

```json
{
  "is_test_net": true,
  "api_key": "",
  "api_secret_key": "",
  "order": {
    "symbol": "BTCUSDT",
    "type": "Market",
    "price": 0, // If type is "Limit" set it as an int greater than 0
    "side": "Buy",
    "qty": 0.014,
    "tp": "0", // see below
    "sl": "{{high}}" // see below
  }
}
```

### TP and SL

You need to set `tp` and `sl` as a string.

- `{{high}}` is an embedded value of TradingView, Also you can set any other TradingView's embedded values.   
- Other methods, you can set as a percent like `"tp": "10%"`. The formula is `(price * percent * 0.1 * qty) + price`. Note that this is calculated from the closing price one minute ago.
- `"tp": "+40", "sl": "-20"` means, `TP: current price + 40` and `SL: current price - 20`.

see [tv.go](pkg/domain/tv.go) or [Bybit API Documentation](https://bybit-exchange.github.io/docs/linear/#:~:text=Transaction%20timestamp-,order,-How%20to%20Subscribe).

Path
-----

| Path     | Method | Description         |
|----------|--------|---------------------|
| /tv      | POST   | Create order        |
| /setting | PUT    | Update your setting |
| /setting | GET    | Get your setting    |

### PUT /setting

You can set the time of day creating order.

Default is all time.

#### Request body

```json
{
  "api_key": "",
  "api_secret_key": "",
  "start_time": "09:00",
  "stop_time": "23:00"
}
```

### GET /setting

Get your setting.

#### Request body

```json
{
  "api_key": "",
  "api_secret_key": ""
}
```

Setup
-----

You can change listen port with `PORT` environment variable.

### Docker

```shell
$ docker run ghcr.io/rluisr/tvbit-bot:latest --name tvbit-bot -p 8080:8080 -d
```

### Binary

1. Download a binary from [Release page.](https://github.com/rluisr/tvbit-bot/releases)
2. `$ ./app`

###  MySQL

tvbit-bot requires MySQL for storing user setting.

Set these environment variables:
- MYSQL_HOST_RW
- MYSQL_HOST_RO
- MYSQL_USER
- MYSQL_PASS
- MYSQL_DB_NAME

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
