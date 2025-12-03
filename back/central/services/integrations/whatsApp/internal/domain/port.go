package domain

import "context"

type IWhatsApp interface {
	SendMessage(ctx context.Context, phoneNumberID uint, msg TemplateMessage) (string, error)
}
