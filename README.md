tvbit-bot
============

[![lint](https://github.com/rluisr/tvbit-bot/actions/workflows/lint.yml/badge.svg?branch=master)](https://github.com/rluisr/tvbit-bot/actions/workflows/lint.yml)

tvbit-bot is TradingView webhook handler for Bybit.

tvbit = T(rading)V(iew) (By)bit

Twitter [@rarirureluis](https://twitter.com/rarirureluis)

Introduction
-------------

1. First of all, register you account. See `PUT /setting`
2. Enable Webhook `https://<domain>/tv`
3. Set an alert with webhook and a message as JSON like below:

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
- Other methods, you can set as a percent like `"tp": "10%"` calculate from mark price.
- `"tp": "+40", "sl": "-20"` means, `TP: mark price + 40` and `SL: mark price - 20`.

see [tv.go](pkg/domain/tv.go)
or [Bybit API Documentation](https://bybit-exchange.github.io/docs/linear/#:~:text=Transaction%20timestamp-,order,-How%20to%20Subscribe)
.

Path
-----

| Path     | Method | Description             |
|----------|--------|-------------------------|
| /tv      | POST   | Create order            |
| /setting | PUT    | Update your setting     |
| /setting | GET    | Get your setting        |
| /wallet  | GET    | Get your wallet balance |

### PUT /setting

Register your account.

Why you have to register an account?  
tvbit-bot fetch your wallet balance and save it to DB. [cron.go](pkg/external/cron.go)

#### Request body

```json
{
  "is_testnet": true,
  "cex": "bybit", 
  "api_key": "",
  "api_secret_key": "",
  "max_position": 0,
  "start_time": "09:00",
  "stop_time": "23:00"
}
```

- **`cex` must be lowercase.**
- `max_position` is how many active position you have
- `start_time` `end_time` is optional.

### [WIP] GET /setting

Get your setting.

#### Response

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

`ghcr.io/rluisr/tvbit-bot:latest`

### Binary

[Release page.](https://github.com/rluisr/tvbit-bot/releases)

### MySQL

tvbit-bot requires MySQL for storing user setting and wallet balance histories.

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

TODO
-----

**There are no plans to support bybit/options/spot trading and other CEX.**

- [x] [bybit/core] Futures
- [x] [bybit/core] USDC
- [x] [bybit/wallet] Deriv
- [x] [bybit/wallet] USDC
- [ ] [core/binance]
