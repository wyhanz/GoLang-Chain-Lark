package llama

import "net/http"

type InferReq struct {
	Msg       string
	Maxtokens int
	Client    *http.Client
}

type TgiPayload struct {
	Inputs string         `json:"inputs"`
	Paras  map[string]int `json:"parameters"`
}

type InferResTgi struct {
	Answer string `json:"generated_text"`
}
