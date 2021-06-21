package model

type Message struct {
	From int `json:"from"`
	To   int `json:"to"`
	Data string `json:"data"`
}
