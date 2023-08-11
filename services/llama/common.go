package llama

import "net/http"

type InferReq struct {
	Msg       string
	Maxtokens int
	Client    *http.Client
}

type TgiPayload struct {
	Inputs  string `json:"prompt"`
	Content string `json:"context"`
	// Paras  map[string]int `json:"parameters"`
}

type InferResTgi struct {
	Answer string `json:"generated_text"`
}

type InferResCode struct {
	Answers []map[string]interface{} `json:"results"`
}
