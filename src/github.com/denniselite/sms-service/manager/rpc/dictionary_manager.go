package rpc

import (
	"encoding/json"
	"github.com/denniselite/go-toolkit/api"
	"github.com/denniselite/go-toolkit/conn"
	"github.com/satori/go.uuid"
	"github.com/denniselite/sms-service/structs"
)

const (
	queueGetCountries = "r.dictionary.GetCountries.v1"
)

type DictionaryManager struct {
	Rmq conn.RmqInt
}

type (
	// Запрос на создание пользователя
	GetCountriesMessage struct {
		Token  string
		Client api.Client
		Body   GetCountriesMessageBody
	}

	GetCountriesMessageBody struct {}

	// ответ на запрос списка стран
	GetCountriesResponse struct {
		Items []structs.Country `json:"items"`
	}
)

// Message: Создание счета в МТ
func (m *DictionaryManager) GetCountries(token string, client api.Client) (countries []structs.Country, err error) {
	msg := GetCountriesMessage{
		Token:  token,
		Client: client,
		Body: GetCountriesMessageBody{},
	}
	body, err := json.Marshal(&msg)
	if err != nil {
		return
	}

	res, err := m.Rmq.Rpc(queueGetCountries, uuid.NewV4().String(), body)
	if err != nil{
		return
	}
	r := new(GetCountriesResponse)
	err = json.Unmarshal(res, r)
	countries = r.Items
	return
}
