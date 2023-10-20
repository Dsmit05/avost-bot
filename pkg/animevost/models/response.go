package models

//go:generate easyjson -all
type Resp struct {
	State State      `json:"state"`
	Items AnimeSpecs `json:"data"`
}

// State содержит дополнительную информацию об ответе
type State struct {
	Status string `json:"status"`
	Rek    int    `json:"rek"`
	Page   int    `json:"page"`
	Count  int    `json:"count"`
}

// ErrorResp ошибочный ответ от api
type ErrorResp struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}
