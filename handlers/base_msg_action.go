package handlers

import (
	"context"
	"fmt"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"gongsheng.cn/agent/services/larkservice"
	"gongsheng.cn/agent/services/llama"
)

type MsgInfo struct {
	msgType   string
	msgId     *string
	chatId    *string
	prompt    string
	fileKey   string
	imageKey  string
	sessionId *string
	mention   []*larkim.MentionEvent
}

type ActionInfo struct {
	handler *MessageHandler
	ctx     *context.Context
	info    *MsgInfo
}

type Action interface {
	Execute(a *ActionInfo) bool
}

type ProcessedUniqueAction struct { //消息唯一性
	// 飞书通过HTTP POST发送Json格式时间数据到用户服务器中。
	// 用户服务器需要在1s内以HTTP 200状态码相应该请求（不需要返回什么Json数据什么的），
	// 否则视为此次事件推送失败，并以5s、5m、1h、6h的间隔重新推送事件，最多重试4次。
}

func (pu *ProcessedUniqueAction) Execute(a *ActionInfo) bool {
	if a.handler.messageCached.IfProcessed(*a.info.msgId) {
		return false
	}
	a.handler.messageCached.TagProcessed(*a.info.msgId)
	return true
}

type ProcessExternalFile struct {
}

func (pf *ProcessExternalFile) Execute(a *ActionInfo) bool {
	if a.info.msgType == "file" {
		replyMsg(*a.ctx, "检测到你发送了一个文件, 回复你自己的文件, 机器人🤖️会根据文件内容以及你说的话给出回答~", a.info.msgId)
		return false
	}
	return true
}

type ProcessRootMsg struct {
}

func (pr *ProcessRootMsg) Execute(a *ActionInfo) bool {
	if *a.info.sessionId != *a.info.msgId {
		content, _ := larkservice.GetLarkClientMsg(*a.info.sessionId)
		fileKey := parseFileKey(content)
		if fileKey != "" {
			a.info.fileKey = fileKey
			a.info.msgType = "post"
			return true
		}
	}
	return true
}

type EasyPrompt struct {
}

func (ep *EasyPrompt) Execute(a *ActionInfo) bool {
	if a.info.msgType == "text" {
		a.info.prompt, _ = llama.BuildPrompt(a.info.prompt)
	} else {
		if a.info.fileKey != "" {
			// 取根节点的文件
			binaryTxT, err := larkservice.GetLarkClientFile(a.info.fileKey, *a.info.sessionId)
			if err != nil {
				fmt.Println(err)
				return false
			}
			fileContent := string(binaryTxT)
			a.info.prompt, _ = llama.BuileFilePrompt(a.info.prompt, fileContent)
		}
	}
	return true
}

type EasyInfer struct {
}

func (ei *EasyInfer) Execute(a *ActionInfo) bool {
	fmt.Println(a.info.prompt)
	inferRes := a.handler.generatedAgent.InferTgi(a.info.prompt, "http://157.148.7.64:38880/generate")
	fmt.Println(inferRes.Answer)
	replyMsg(*a.ctx, inferRes.Answer, a.info.msgId)
	return true
}
