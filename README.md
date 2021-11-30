# WeChat SDK for Go
![Go](https://github.com/amazing-gao/wechat/workflows/Go/badge.svg?branch=release-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/amazing-gao/wechat)](https://goreportcard.com/report/github.com/amazing-gao/wechat)
[![pkg](https://img.shields.io/badge/dev-reference-007d9c?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/amazing-gao/wechat/v2?tab=doc)
![version](https://img.shields.io/badge/version-v2-green)

使用Golang开发的微信SDK，简单、易用。
>注意：当前版本为v2版本，v1版本已废弃


## 文档 && 例子
[API列表](https://github.com/amazing-gao/wechat/tree/v2/doc/api)

## 快速开始
```
import "github.com/amazing-gao/wechat/v2"
```

以下是一个微信公众号处理消息接收以及回复的例子：

```go
//使用memcache保存access_token，也可选择redis或自定义cache
wc := wechat.NewWechat()
memory := cache.NewMemory()
cfg := &offConfig.Config{
    Server:    "https://api.weixin.qq.com",
    AppID:     "xxx",
    AppSecret: "xxx",
    Token:     "xxx",
    EncodingAESKey: "xxxx",
    Cache: memory,
}
officialAccount := wc.GetOfficialAccount(cfg)

// 传入request和responseWriter
server := officialAccount.GetServer(req, rw)
//设置接收消息的处理方法
server.SetMessageHandler(func(msg *message.MixMessage) *message.Reply {

    //回复消息：演示回复用户发送的消息
    text := message.NewText(msg.Content)
    return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
})

//处理消息接收以及回复
err := server.Serve()
if err != nil {
    fmt.Println(err)
    return
}
//发送回复的消息
server.Send()

```

## 目录说明
- officialaccount: 微信公众号API
- miniprogram: 小程序API
- minigame:小游戏API
- pay:微信支付API
- openplatform:开放平台API
- work:企业微信
- aispeech:智能对话
- doc: api文档

## License

Apache License, Version 2.0
