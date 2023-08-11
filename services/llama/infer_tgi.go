package llama

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"gongsheng.cn/agent/global"
)

func constructPromptChain(prompts []string) string {
	finalPrompt := ""
	for _, p := range prompts {
		finalPrompt += p
	}
	return finalPrompt
}

func BuildPrompt(msg string) (string, error) {
	//目前先做单轮prompt构成
	prompt := msg + "[SPLIT]"

	return prompt, nil
}

func BuileFilePrompt(msg, fileContent string) (string, error) {
	prompts := []string{
		msg,
		"[SPLIT]",
		fileContent,
	}
	prompt := constructPromptChain(prompts)
	return prompt, nil
}

func promptSpliter(msg string) (string, string, error) {
	parts := strings.Split(msg, "[SPLIT]")

	return parts[0], parts[1], nil
}

func (ir *InferReq) InferTgi(msg, url string) *InferResTgi {
	timeOutDuration := time.Duration(global.INFER_TIME_REQ) * time.Second
	ir.Client = &http.Client{Timeout: timeOutDuration}
	prompt, content, _ := promptSpliter(msg)
	requestBody := TgiPayload{
		Inputs:  prompt,
		Content: content,
		// Paras: map[string]int{
		// 	"max_new_tokens": 256,
		// },
	}
	responseBody := &InferResCode{}
	returnRes := &InferResTgi{}
	err := ir.sendReqWithRetry(url, requestBody, responseBody, ir.Client, 3)
	if err != nil {
		fmt.Printf("error when sending req, %s", err)
	}
	returnRes.Answer = responseBody.Answers[0]["text"].(string)
	return returnRes
}

func (ir *InferReq) sendReqWithRetry(url string,
	requestBody, responseBody interface{},
	client *http.Client, maxRetries int) error {

	var requestBodyData []byte
	var err error
	requestBodyData, err = json.Marshal(requestBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewReader(requestBodyData))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	var response *http.Response
	var retry int

	for retry = 0; retry <= maxRetries; retry++ {
		if retry > 0 {
			req.Body = io.NopCloser(bytes.NewReader(requestBodyData))
		}

		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

		response, err = client.Do(req)
		if err != nil || response.StatusCode < 200 || response.StatusCode >= 300 {
			fmt.Println(err)
			body, _ := io.ReadAll(response.Body)
			fmt.Println("body", string(body))

			if retry == maxRetries {
				break
			}
			time.Sleep(time.Duration(retry+1) * time.Second)
		} else {
			break
		}
	}
	if response != nil {
		defer response.Body.Close()
	}
	if response == nil || response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("%s api failed after %d retries", strings.ToUpper("infer"), retry)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, responseBody)
	if err != nil {
		return err
	}

	return nil
}
