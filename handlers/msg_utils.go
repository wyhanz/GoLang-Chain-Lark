package handlers

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
)

func msgFilter(msg string) string {
	//replace @到下一个非空的字段 为 ''
	regex := regexp.MustCompile(`@[^ ]*`)
	return regex.ReplaceAllString(msg, "")

}

func parseContent(content string) string {
	//"{\"text\":\"@_user_1  hahaha\"}",
	//only get text content hahaha

	var contentMap map[string]interface{}
	err := json.Unmarshal([]byte(content), &contentMap)
	if err != nil {
		fmt.Println(err)
	}
	if contentMap["text"] == nil {
		return ""
	}
	text := contentMap["text"].(string)
	return msgFilter(text)
}

func parseFileKey(content string) string {
	//"{\"text\":\"@_user_1  hahaha\"}",
	//only get text content hahaha

	var contentMap map[string]interface{}
	err := json.Unmarshal([]byte(content), &contentMap)
	if err != nil {
		fmt.Println(err)
	}
	if contentMap["file_key"] == nil {
		return ""
	}
	if filepath.Ext(contentMap["file_name"].(string)) != ".txt" {
		return ""
	}
	fileKey := contentMap["file_key"].(string)
	return fileKey
}
