package lib

import (
	"WarnNotify/internal"
	"WarnNotify/types"
	"context"
	"time"
)

var (
	_ Notify = &internal.UnsupportedNotify{}
	_ Notify = &internal.TelegramNotify{}
)

type Notify interface {
	WarnMessage(text string) error
	WarnStructMessage(message types.StructMessage) error
	Watch(getMeta func() string, duration time.Duration) (context.CancelFunc, error)
}

func NewNotify(config Config) Notify {
	switch config.NotifyChannel {
	case types.NotifyChannelTypeTelegram:
		return internal.NewTelegram(config.TelegramBotToken, config.TelegramNotifyUsers, config.NeedPing)
	default:
		panic("unsupported notify channel")
		return &internal.UnsupportedNotify{}
	}
}
