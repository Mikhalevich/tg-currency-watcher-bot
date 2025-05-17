package messagesender

import (
	"context"
)

type Noop struct {
}

func NewNoop() *Noop {
	return &Noop{}
}

func (n *Noop) SendTextMessage(ctx context.Context, chatID int64, text string) {
}
