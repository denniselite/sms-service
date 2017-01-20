package structs

type Transport interface {
	SetOptions(options map[string]string)
	SendSms(phone string, message string) error
	GetName() string
}
