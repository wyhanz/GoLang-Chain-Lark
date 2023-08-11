package llama

// func (ir *InferReq) InferTgi(msg, url string) *InferResTgi {
// 	timeOutDuration := time.Duration(global.INFER_TIME_REQ) * time.Second
// 	ir.Client = &http.Client{Timeout: timeOutDuration}
// 	requestBody := TgiPayload{
// 		Inputs: msg,
// 		Paras: map[string]int{
// 			"max_new_tokens": 256,
// 		},
// 	}
// 	responseBody := &InferResTgi{}
// 	err := ir.sendReqWithRetry(url, requestBody, responseBody, ir.Client, 3)
// 	if err != nil {
// 		fmt.Printf("error when sending req, %s", err)
// 	}
// 	return responseBody
// }

// func (ir *InferReq) sendReqWithRetry(url string,
// 	requestBody, responseBody interface{},
// 	client *http.Client, maxRetries int) error {

// 	var requestBodyData []byte
// 	var err error
// 	requestBodyData, err = json.Marshal(requestBody)
// 	if err != nil {
// 		return err
// 	}
// 	req, err := http.NewRequest("POST", url, bytes.NewReader(requestBodyData))
// 	if err != nil {
// 		return err
// 	}
// 	req.Header.Set("Content-Type", "application/json")

// 	var response *http.Response
// 	var retry int

// 	for retry = 0; retry <= maxRetries; retry++ {
// 		if retry > 0 {
// 			req.Body = io.NopCloser(bytes.NewReader(requestBodyData))
// 		}
// 		response, err = client.Do(req)
// 		if err != nil || response.StatusCode < 200 || response.StatusCode >= 300 {

// 			body, _ := io.ReadAll(response.Body)
// 			fmt.Println("body", string(body))

// 			if retry == maxRetries {
// 				break
// 			}
// 			time.Sleep(time.Duration(retry+1) * time.Second)
// 		} else {
// 			break
// 		}
// 	}
// 	if response != nil {
// 		defer response.Body.Close()
// 	}
// 	if response == nil || response.StatusCode < 200 || response.StatusCode >= 300 {
// 		return fmt.Errorf("%s api failed after %d retries", strings.ToUpper("infer"), retry)
// 	}

// 	body, err := io.ReadAll(response.Body)
// 	if err != nil {
// 		return err
// 	}
// 	err = json.Unmarshal(body, responseBody)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
