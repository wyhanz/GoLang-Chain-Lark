package handlers

import (
	"context"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"gongsheng.cn/agent/initialization"
	"gongsheng.cn/agent/logger"
	"gongsheng.cn/agent/services/llama"
)

type MessageHandlerInterface interface {
	msgReceivedHandler(ctx context.Context, event *larkim.P2MessageReceiveV1) error
	// cardHandler(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error)
}

type HandlerType string

// handlers 所有消息类型类型的处理器
var handlers MessageHandlerInterface

func InitHandlers(gpt llama.InferReq, config initialization.Config) {
	handlers = NewMessageHandler(gpt, config)
}

func Handler(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	return handlers.msgReceivedHandler(ctx, event)
}

func ReadHandler(ctx context.Context, event *larkim.P2MessageReadV1) error {
	readerId := event.Event.Reader.ReaderId.OpenId
	//fmt.Printf("msg is read by : %v \n", *readerId)
	logger.Debugf("msg is read by : %v \n", *readerId)

	return nil
}
