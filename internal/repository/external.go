package repository

import (
	"context"

	"github.com/hanifbg/landing_backend/internal/model/response"
)

type WhatsApp interface {
	SendMessage(phoneNumber, message string) (err error)
}

type TelegramAPI interface {
	SendMessage(ctx context.Context, ChatID int64, ThreadID int64, message string) (resp *response.Message, err error)
}
