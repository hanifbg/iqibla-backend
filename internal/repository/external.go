package repository

type WhatsApp interface {
	SendMessage(phoneNumber, message string) (err error)
}
