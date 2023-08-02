package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	sdkginext "github.com/larksuite/oapi-sdk-gin"
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"gongsheng.cn/agent/handlers"
	"gongsheng.cn/agent/initialization"
	"gongsheng.cn/agent/logger"
	"gongsheng.cn/agent/services/llama"
	"gongsheng.cn/agent/utils"
)

func main() {
	gpt := &llama.InferReq{Msg: "msg", Maxtokens: 256}

	config := initialization.GetConfig()
	initialization.LoadLarkClient(*config)
	handlers.InitHandlers(*gpt, *config)
	eventHandler := dispatcher.NewEventDispatcher(config.FeishuVerifiedToken, config.FeishuEncryptKey).
		OnP2MessageReceiveV1(handlers.Handler).
		OnP2MessageReadV1(func(ctx context.Context, event *larkim.P2MessageReadV1) error {
			logger.Debugf("收到请求 %v", event.RequestURI)
			return handlers.ReadHandler(ctx, event)
		})

	// 创建卡片行为处理器
	cardHandler := larkcard.NewCardActionHandler(config.FeishuAppId, config.FeishuAppSecret, func(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error) {
		fmt.Println(larkcore.Prettify(cardAction))

		// 返回卡片消息
		//return getCard(), nil

		//custom resp
		//return getCustomResp(),nil

		// 无返回值
		return nil, nil
	})
	g := gin.Default()
	g.GET("ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	g.POST("/", func(ctx *gin.Context) {
		var challenge utils.ChallengeReq
		ctx.ShouldBindJSON(&challenge)
		ctx.JSON(200, gin.H{
			"challenge": challenge.Challenge,
		})
	})
	g.POST("/webhook/event", sdkginext.NewEventHandlerFunc(eventHandler))
	g.POST("/webhook/card", sdkginext.NewCardActionHandlerFunc(cardHandler))

	// 启动服务
	err := g.Run(":9999")
	if err != nil {
		panic(err)
	}
}
