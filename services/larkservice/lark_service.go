package larkservice

import (
	"context"
	"encoding/json"
	"io"

	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
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

func GetLarkClientMsg(msgId string) (string, error) {

	var err error
	client := initialization.GetLarkClient()

	req := larkim.NewGetMessageReqBuilder().
		MessageId(msgId).
		Build()

	resp, err := client.Im.Message.Get(context.Background(), req)
	if err != nil {
		return "", nil
	}
	var respMap map[string]interface{}
	// var ContentMap map[string]interface{}
	err = json.Unmarshal(resp.RawBody, &respMap)
	if err != nil {
		return "", nil
	}

	//这里...
	res := respMap["data"].(map[string]interface{})["items"].([]interface{})[0].(map[string]interface{})["body"].(map[string]interface{})["content"].(string)
	return res, err
}
