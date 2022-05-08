tvbit-bot
============

[![lint](https://github.com/rluisr/tvbit-bot/actions/workflows/lint.yml/badge.svg?branch=master)](https://github.com/rluisr/tvbit-bot/actions/workflows/lint.yml)

[English README](README_en.md)

tvbit-bot は TradingView のアラートから Bybit へ注文を行う BOT です。

tvbit = T(rading)V(iew) (By)bit

Introduction
-------------

1. Webhook の送信先を設定する `https://<domain>/tv`
2. 以下の JSON を参考に TradingView のアラートメッセージを設定しアラートを作成する

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

`{{high}}` is an embedded value of TradingView, Also you can set any other TradingView's embedded values.   
Other methods, you can set as a percent like `"tp": "10%"`. It means if price is 30,000 and qty is 0.1, TP is set `30,033.5`

[tv.go](pkg/domain/tv.go) または [Bybit API Documentation](https://bybit-exchange.github.io/docs/linear/#:~:text=Transaction%20timestamp-,order,-How%20to%20Subscribe) を参照してください。

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

**All times are calculated in UTC.**

#### Request body

```json
{
  "api_key": "",
  "api_secret_key": "",
  "start_time": "09:00",
  "stop_time": "23:00"
}
```

#### Response

nothing

### GET /setting

Get your setting.

#### Request body

```json
{
  "api_key": "",
  "api_secret_key": ""
}
```

#### Response

```json
{
  "api_key": "",
  "api_secret_key": "",
  "start_time": "09:00",
  "stop_time": "23:00"
}
```

Setup
-----

環境変数 `PORT` でリッスンポートを上書きできます。

### Docker

```shell
$ docker run ghcr.io/rluisr/tvbit-bot:latest --name tvbit-bot -p 8080:8080 -d
```

### Binary

1. [Release](https://github.com/rluisr/tvbit-bot/releases) からバイナリをダウンロード
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

私はどなたでも利用できるようにこのアプリケーションを公開しています。  
しかし私は裏切る可能性もあるのでテストネットでお試しください。  
もちろん本番環境でも利用できますが私は一切の責任を負いません。

Powered by [HCloud Ltd](https://hcloud.ltd)

Limitation
----------

tvbit-bot は現在 close/cancel ポジションをサポートしていません。

Welcome your PR.

Twitter [@rarirureluis](https://twitter.com/rarirureluis)
