package handlers

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"gongsheng.cn/agent/initialization"
	"gongsheng.cn/agent/utils"
)

func replyMsg(ctx context.Context, msg string, msgId *string) error {
	client := initialization.GetLarkClient()
	msg = utils.EscapeJsonChars(msg)
	content := larkim.NewTextMsgBuilder().
		Text(msg).
		Build()

	resp, err := client.Im.Message.Reply(ctx, larkim.NewReplyMessageReqBuilder().
		MessageId(*msgId).
		Body(larkim.NewReplyMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeText).
			Uuid(uuid.New().String()).
			Content(content).
			Build()).
		Build())

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return errors.New(resp.Msg)
	}
	// 业务处理
	fmt.Println(larkcore.Prettify(resp))
	return nil
}
