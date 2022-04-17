tvbit-bot
============

[![lint](https://github.com/rluisr/tvbit-bot/actions/workflows/lint.yml/badge.svg?branch=master)](https://github.com/rluisr/tvbit-bot/actions/workflows/lint.yml)

[English README](README_en.md)

tvbit-bot は TradingView のアラートから Bybit へ注文を行う BOT です。

tvbit = T(rading)V(iew) (By)bit

Introduction
-------------

1. 以下の JSON を参考に TradingView のアラートメッセージを設定しアラートを作成する

```json
{
  "is_test_net": true,
  "api_key": "",
  "api_secret_key": "",
  "order": {
    "symbol": "BTCUSDT",
    "type": "Market",
    "price": 0,
    "side": "Sell",
    "qty": 0.028,
    "tp": 0,
    "sl": {{high}}
  }
}
```

[tv.go](pkg/domain/tv.go) または [Bybit API Documentation](https://bybit-exchange.github.io/docs/linear/#:~:text=Transaction%20timestamp-,order,-How%20to%20Subscribe) を参照してください。

tvbit-bot は Bybit API に準拠しているつもりです。

Setup
-----

環境変数 `PORT` でリッスンポートを上書きできます。

### Docker

```shell
$ docker run ghcr.io/rluisr/tvbit-bot:latest --name tvbit-bot -p 8080:8080 -d
```

### Binary

1. [release page.](https://github.com/rluisr/tvbit-bot/releases) からバイナリをダウンロード
2. `$ ./app`

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
