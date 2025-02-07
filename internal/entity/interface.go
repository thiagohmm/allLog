package entity

import "context"

type MessageRepository interface {
	SaveMessage(ctx context.Context, message Message) error
}
