package telegram

import "context"

type Notify struct {
	token string
	users []int64

	queue  chan *request
	cancel context.CancelFunc
}

type request struct {
	User int64  `json:"chat_id"`
	Msg  string `json:"text"`
}
