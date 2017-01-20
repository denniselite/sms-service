package components

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type ThaiBulkSms struct {
	Options map[string]string
}

func (t *ThaiBulkSms) SetOptions(options map[string]string) {
	t.Options = options
}

func (t *ThaiBulkSms) GetName() (name string) {
	return "ThaiBulkSms"
}

func (t *ThaiBulkSms) CorrectPhone(phone string) string {
	if !strings.HasPrefix(phone, "0") {
		phone = "0" + phone
	}

	return phone
}

func (t *ThaiBulkSms) SendSms(phone string, message string) error {
	v := url.Values{}
	v.Set("username", t.Options["login"])
	v.Set("password", t.Options["password"])
	v.Set("msisdn", t.CorrectPhone(phone))
	v.Set("message", message)

	rb := *strings.NewReader(v.Encode())

	client := &http.Client{}
	req, err := http.NewRequest("POST", t.Options["url"], &rb)
	if err != nil {
		return err
	}

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

	var response ThaiBulkSmsResponse
	err = xml.Unmarshal(contents, &response)
	if err != nil {
		return err
	}

	if response.Status != 1 {
		return errors.New(fmt.Sprintf("Error: %s", response.Detail))
	}

	return nil
}

type ThaiBulkSmsResponse struct {
	Status int    `xml:"Status"`
	Detail string `xml:"Detail"`
}
