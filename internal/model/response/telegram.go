package response

import "github.com/go-telegram/bot/models"

type Message struct {
	Ok      bool           `json:"ok"`
	Message models.Message `json:"message"`
}
