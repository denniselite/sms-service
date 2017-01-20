package manager

import (
	"github.com/denniselite/go-toolkit/conn"
	. "github.com/denniselite/go-toolkit/errors"
	"github.com/denniselite/go-toolkit/api"
	"github.com/denniselite/sms-service/structs"
	"github.com/denniselite/sms-service/libs"
	"github.com/denniselite/sms-service/manager/rpc"
	"log"
)

type SmsManager struct {
	Rmq          conn.RmqInt
	Db           *conn.DB
	transports   map[string][]structs.Transport
	typeRegistry map[string]structs.Transport
	SmsProviders []structs.SmsProvider
	DictionaryManager         rpc.DictionaryManager
}

// Инициализация менеджера
func (m *SmsManager) Run(rmq conn.RmqInt) {
	m.Rmq = rmq
	m.DictionaryManager.Rmq = rmq
	// Подготавливаем набор провайдеров
	OopsT(api.SysToken, m.initTransports())

	go func() {
		OopsT(api.SysToken, m.Rmq.ConsumeRpc("r.sms.SendSms.v1", m.SendSms, -1))
	}()
}

func (m *SmsManager) addTransport(transport components.Transport) {
	m.typeRegistry[transport.GetName()] = transport
}

// Подготовка набора провайдеров для телефонных кодов стран
func (m *SmsManager) initTransports() (err error) {

	// Регистриуем провайдеры
	m.typeRegistry = make(map[string]structs.Transport)

	m.addTransport(&components.SmsDirect{})
	m.addTransport(&components.FiboSms{})
	m.addTransport(&components.ISMSIndonesia{})
	m.addTransport(&components.ThaiBulkSms{})

	// Устанавливаем опции провайдеров
	for _, provider := range m.SmsProviders {
		m.typeRegistry[provider.Transport].SetOptions(provider.Settings)
	}

	m.transports = make(map[string][]structs.Transport)

	// Достаем список телефонных кодов стран
	countries, err := m.DictionaryManager.GetCountries(api.GenerateToken(), api.Client{})
	if err != nil {
		return
	}

	// Сначала распределяем тех, у кого есть соответствие по коду страны
	for _, provider := range m.SmsProviders {
		if (len(provider.Codes)) > 0 {
			for _, providerCode := range provider.Codes {
				m.transports[providerCode] = append(m.transports[providerCode], m.typeRegistry[provider.Transport])
			}
		}
	}

	// Дополняем набор провайдеров для кодов стран провайдерами по умолчанию
	for _, country := range countries {
		for _, provider := range m.SmsProviders {
			if (len(provider.Codes)) == 0 {
				m.transports[country.CodePhone] = append(m.transports[country.CodePhone], m.typeRegistry[provider.Transport])
			}
		}
	}

	log.Println(api.SysToken, "SMS transports initialized")

	return
}