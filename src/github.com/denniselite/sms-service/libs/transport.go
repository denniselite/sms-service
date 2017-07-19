package components

// Интерфейс отправки СМС
type Transport interface {
	SetOptions(options map[string]string)
	SendSms(phone string, message string) error
	GetName() string
}
