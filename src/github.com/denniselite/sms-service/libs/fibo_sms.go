package components

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type FiboSms struct {
	Options map[string]string
}

func (t *FiboSms) SetOptions(options map[string]string) {
	t.Options = options
}

func (t *FiboSms) GetName() (name string) {
	return "FiboSms"
}

func (t *FiboSms) SendSms(phone string, message string) error {
	req, err := http.NewRequest("GET", t.Options["url"], nil)
	if err != nil {
		return err
	}

	v := req.URL.Query()
	v.Add("clientNo", t.Options["login"])
	v.Add("clientPass", t.Options["password"])
	v.Add("senderName", t.Options["from"])
	v.Add("phoneNumber", phone)
	v.Add("smsMessage", message)
	v.Add("smsGUID", "1234")
	v.Add("serviceType", t.Options["serviceType"])
	req.URL.RawQuery = v.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("SMS sending error. HTTP Status Code: %d", resp.StatusCode))
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var resultString string
	var response FiboResponse
	if err := xml.Unmarshal(contents, &resultString); err != nil {
		return err
	}
	if err := xml.Unmarshal([]byte(resultString), &response); err != nil {
		return err
	}
	if response.Code != 0 {
		return errors.New(fmt.Sprintf("Error: %s", response.Message))
	}

	return nil
}

type FiboResponse struct {
	Code    int
	Message string
	Time    string
}
