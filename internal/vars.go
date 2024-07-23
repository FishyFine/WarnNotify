package internal

import (
	"WarnNotify/types"
	"context"
	"time"
)

type UnsupportedNotify struct{}

func (_ *UnsupportedNotify) WarnMessage(text string) error {
	//TODO implement me
	panic("implement me")
}

func (_ *UnsupportedNotify) WarnStructMessage(message types.StructMessage) error {
	//TODO implement me
	panic("implement me")
}

func (_ *UnsupportedNotify) Watch(getMeta func() string, duration time.Duration) (context.CancelFunc, error) {
	//TODO implement me
	panic("implement me")
}
