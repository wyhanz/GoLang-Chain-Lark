package larkservice

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	larkdocx "github.com/larksuite/oapi-sdk-go/v3/service/docx/v1"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	larkwiki "github.com/larksuite/oapi-sdk-go/v3/service/wiki/v2"
	"gongsheng.cn/agent/initialization"
)

func GetLarkClientFile(fileKey, msgId string) ([]byte, error) {
	client := initialization.GetLarkClient()
	// 创建请求对象
	req := larkim.NewGetMessageResourceReqBuilder().
		MessageId(msgId).
		FileKey(fileKey).
		Type("file").
		Build()

	resp, err := client.Im.MessageResource.Get(context.Background(), req)
	if err != nil {
		return nil, err
	}

	binaryFile, err := io.ReadAll(resp.File)
	if err != nil {
		return nil, err
	}

	return binaryFile, nil
}

func GetLarkClientMsg(msgId string) (map[string]interface{}, error) {

	var err error
	client := initialization.GetLarkClient()

	req := larkim.NewGetMessageReqBuilder().
		MessageId(msgId).
		Build()

	resp, err := client.Im.Message.Get(context.Background(), req)
	if err != nil {
		return nil, nil
	}
	var respMap map[string]interface{}
	err = json.Unmarshal(resp.RawBody, &respMap)
	if err != nil {
		return nil, nil
	}

	res := respMap["data"].(map[string]interface{})["items"].([]interface{})[0].(map[string]interface{})
	return res, err
}

func GetLarkDocsContent(docId string) (string, error) {
	client := initialization.GetLarkClient()

	req := larkdocx.NewRawContentDocumentReqBuilder().
		DocumentId(docId).
		Lang(1).
		Build()

	resp, err := client.Docx.Document.RawContent(context.Background(), req)

	// 处理错误
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return "", fmt.Errorf("Fail; Code: %s, Msg: %s", fmt.Sprint(resp.Code), resp.Msg)
	}
	var respMap map[string]interface{}
	err = json.Unmarshal(resp.RawBody, &respMap)
	if err != nil {
		return "", err
	}
	// code :=
	res := respMap["data"].(map[string]interface{})["content"].(string)
	return res, nil
}

func GetLarkWikiContent(wikiId string) (string, error) {
	var (
		nodeId  string
		content string
	)
	var err error
	if nodeId, err = GetLarkWikiInfo(wikiId); err != nil {
		return "", err
	}
	if content, err = GetLarkDocsContent(nodeId); err != nil {
		return "", err
	}

	return content, err
}

func GetLarkWikiInfo(wikiId string) (string, error) {
	client := initialization.GetLarkClient()

	req := larkwiki.NewGetNodeSpaceReqBuilder().
		Token(wikiId).
		Build()

	resp, err := client.Wiki.Space.GetNode(context.Background(), req)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// 服务端错误处理
	if !resp.Success() {
		fmt.Println(resp.Code, resp.Msg, resp.RequestId())
		return "", err
	}
	var respMap map[string]interface{}
	err = json.Unmarshal(resp.RawBody, &respMap)
	if err != nil {
		return "", err
	}
	objToken := respMap["data"].(map[string]interface{})["node"].(map[string]interface{})["obj_token"].(string)
	objType := respMap["data"].(map[string]interface{})["node"].(map[string]interface{})["obj_type"].(string)
	if objType != "docx" {
		return "", fmt.Errorf("不太支持除docx外，其他格式文档")
	}
	return objToken, nil
}
