package repository

import "github.com/hanifbg/landing_backend/internal/model/entity"

type Mailer interface {
	Send(from, to, subject, body string) error
	SendOrderConfirmation(order *entity.Order) error
}
