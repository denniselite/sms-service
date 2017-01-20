package manager

import (
	"encoding/json"
	"github.com/denniselite/go-toolkit/api"
	"github.com/denniselite/go-toolkit/conn"
	"gopkg.in/validator.v2"
	"log"
	"errors"
	"fmt"
)

type SendSmsRequest struct {
	Token  string
	Client api.Client
	Body   SendSmsRequestBody
}

type SendSmsRequestBody struct {
	PhoneCode string `validate:"nonzero"`
	PhoneNumber string `validate:"nonzero"`
	Message string `validate:"nonzero"`
}

// Пример RPC метода обработчика
func (m *SmsManager) SendSms(data []byte) (res conn.RmqMessage, err error) {
	rq := new(SendSmsRequest)
	err = json.Unmarshal(data, rq)
	if err != nil {
		return
	}

	err = validator.Validate(rq)
	if err != nil {
		return
	}

	log.Printf("%s Send sms request handled", rq.Token)

	phoneString := rq.Body.PhoneCode + rq.Body.PhoneNumber
	countryTransports := m.transports[rq.Body.PhoneCode]

	// Заранее инициализируем ошибку отсутствия провайдера для выбранного телефонного кода
	err = errors.New(fmt.Sprintf("No transports found for phone country code: %s", rq.Body.PhoneCode))
	for _, transport := range countryTransports {

		// Отправка сообщения
		err = transport.SendSms(phoneString, rq.Body.Message)

		// Если отправка удалась - выходим, иначе - продолжаем пытаться отправить
		// через останых провайдеров, доступных для этого кода страны
		if err == nil {
			log.Printf("%s Send sms via %s to %s; message: %s", rq.Token, transport.GetName(), phoneString, rq.Body.Message)
			break
		}

	}

	if err != nil {
		return
	}

	res = new(conn.EmptyResponse)
	return
}
