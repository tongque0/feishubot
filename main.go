package main

import (
	"context"
	"feishu/callback"
	"feishu/config"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkws "github.com/larksuite/oapi-sdk-go/v3/ws"
)

func main() {
	// 注册事件回调
	eventHandler := dispatcher.NewEventDispatcher("", "").OnP2MessageReceiveV1(callback.MessageReceive)

	// 启动 WebSocket 客户端
	wsClient := larkws.NewClient(config.GetConf().AppID, config.GetConf().AppSecret,
		larkws.WithEventHandler(eventHandler),
		larkws.WithLogLevel(larkcore.LogLevelDebug),
	)

	if err := wsClient.Start(context.Background()); err != nil {
		panic(err)
	}
}
