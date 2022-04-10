tvbit-bot
============

[English](README_en.md)

tvbit-bot は TradingView のアラートから Bybit へ注文を行う BOT です。

tvbit = T(rading)V(iew) (By)bit

Introduction
-------------

1. アプリを動かす
2. 以下の JSON を参考に TradingView のアラートメッセージを設定しアラートを作成する

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

[tv.go](pkg/domain/tv.go) または [Bybit API Documentation](https://bybit-exchange.github.io/docs/linear/#:~:text=Transaction%20timestamp-,order,-How%20to%20Subscribe). を参照してください。

tvbit-bot は Bybit API に準拠しているつもりです。

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