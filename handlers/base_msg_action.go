package handlers

import (
	"context"
	"fmt"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"gongsheng.cn/agent/global"
	"gongsheng.cn/agent/services/larkservice"
	"gongsheng.cn/agent/services/llama"
)

var MsgTypeMapper = map[int]string{
	global.TEXT_LARK: "text",
	global.FILE_LARK: "file",
	global.DOCS_LARK: "lark_docs",
	global.WIKI_LARK: "lark_wiki",
}

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
		content, err := larkservice.GetLarkClientMsg(*a.info.sessionId)
		if err != nil {
			fmt.Println("Get Msg Failed %s", err)
			return true
		}

		if content["content"].(string) == "file" {
			fileKey := parseFileKey(content["content"].(string))
			if fileKey != "" {
				a.info.fileKey = fileKey
				a.info.msgType = "root_file" //处理标识
				return true
			}
		}
	}
	return true
}

type ProcessLarkWiki struct {
}

// Need careful consideration and construction
func (pl *ProcessLarkWiki) Execute(a *ActionInfo) bool {
	// 判断是否是飞书云文档
	fileToken, processedPrompt, LarkMsgNum := judgeIfLarkWiki(a.info.prompt)
	if LarkMsgNum != 0 {
		a.info.msgType = MsgTypeMapper[LarkMsgNum]
		a.info.fileKey = fileToken
		a.info.prompt = processedPrompt
		return true
	}
	return true
}

type EasyPrompt struct {
}

func (ep *EasyPrompt) Execute(a *ActionInfo) bool {
	if a.info.msgType == "text" {
		a.info.prompt, _ = llama.BuildPrompt(a.info.prompt)
	} else if a.info.msgType == "root_file" {
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
	} else if a.info.msgType == "lark_docs" || a.info.msgType == "lark_wiki" {
		var wikiContent string
		var err error
		if a.info.msgType == "lark_docs" {
			wikiContent, err = larkservice.GetLarkDocsContent(a.info.fileKey)
		} else {
			wikiContent, err = larkservice.GetLarkWikiContent(a.info.fileKey)
		}

		if err != nil {
			replyMsg(*a.ctx, fmt.Sprintf("无法获取所提供的wiki, 这可能是因为机器人没有阅读你文档的权限～, %s", err), a.info.msgId)
			return false
		}
		fmt.Println(wikiContent)
		a.info.prompt, _ = llama.BuileFilePrompt(a.info.prompt, wikiContent)
	}
	return true
}

type EasyInfer struct {
}

func (ei *EasyInfer) Execute(a *ActionInfo) bool {
	fmt.Println(a.info.prompt)
	url := fmt.Sprintf("%s/generate", a.handler.config.LlamaUrl)
	inferRes := a.handler.generatedAgent.InferTgi(a.info.prompt, url)
	fmt.Println(inferRes.Answer)
	replyMsg(*a.ctx, inferRes.Answer, a.info.msgId)
	return true
}
