package handlers

import (
	"context"
	"fmt"
	"strings"

	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"gongsheng.cn/agent/initialization"
	"gongsheng.cn/agent/logger"
	"gongsheng.cn/agent/services/llama"
	"gongsheng.cn/agent/utils"
)

// 责任链
func chain(data *ActionInfo, actions ...Action) bool {
	for _, v := range actions {
		if !v.Execute(data) {
			return false
		}
	}
	return true
}

type MessageHandler struct {
	config         initialization.Config
	generatedAgent llama.InferReq
	messageCached  utils.MsgCacheInterface
}

func (m MessageHandler) msgReceivedHandler(ctx context.Context, event *larkim.P2MessageReceiveV1) error {
	logger.Debug("收到消息：", larkcore.Prettify(event.Event.Message))
	fmt.Println(larkcore.Prettify(event.Event.Message))

	msgType, _ := judgeMessageType(event)
	content := event.Event.Message.Content
	msgId := event.Event.Message.MessageId
	rootId := event.Event.Message.RootId
	chatId := event.Event.Message.ChatId
	mention := event.Event.Message.Mentions

	sessionId := rootId
	if sessionId == nil || *sessionId == "" {
		sessionId = msgId
	}

	msgInfo := MsgInfo{
		// handlerType: handlerType,
		msgType:   msgType,
		msgId:     msgId,
		chatId:    chatId,
		prompt:    strings.Trim(parseContent(*content), " "),
		fileKey:   strings.Trim(parseFileKey(*content), " "),
		sessionId: sessionId,
		mention:   mention,
	}

	data := &ActionInfo{
		handler: &m,
		ctx:     &ctx,
		info:    &msgInfo,
	}

	actions := []Action{
		&ProcessedUniqueAction{},
		&ProcessExternalFile{},
		&ProcessRootMsg{},
		&EasyPrompt{},
		&EasyInfer{},
	}
	chain(data, actions...)

	return nil
}

var _ MessageHandlerInterface = (*MessageHandler)(nil) //接口保证

func NewMessageHandler(generatedAgent llama.InferReq, config initialization.Config) MessageHandlerInterface {
	return &MessageHandler{
		config:         config,
		generatedAgent: generatedAgent,
		messageCached:  utils.GetMsgCache(),
	}
}

func judgeMessageType(event *larkim.P2MessageReceiveV1) (string, error) {
	msgType := event.Event.Message.MessageType

	switch *msgType {
	case "text", "post", "file":
		return *msgType, nil
	default:
		return "", fmt.Errorf("unknown message type: %v", *msgType)
	}
}
