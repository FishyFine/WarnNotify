package telegram

import "context"

func NewTelegram(token string, users []int64, needPing bool) *Notify {
	ctx, cancel := context.WithCancel(context.Background())
	n := &Notify{
		token: token,
		users: users,

		queue:  make(chan *request, 1024),
		cancel: cancel,
	}
	if needPing {
		if err := n.WarnMessage("Ping"); err != nil {
			panic(err)
		}
	}
	go func(ctx context.Context) {
		n.asyncRequest()
	}(ctx)
	return n
}
