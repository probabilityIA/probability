package domain

import "context"

type IWhatsappClient interface {
	SendMessage(ctx context.Context, to string, message string) error
}
