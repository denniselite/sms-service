package components

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type SmsDirect struct {
	Options map[string]string
}

func (t *SmsDirect) SetOptions(options map[string]string) {
	t.Options = options
}

func (t *SmsDirect) GetName() (name string) {
	return "SmsDirect"
}

func (t *SmsDirect) SendSms(phone string, message string) error {
	v := url.Values{}
	v.Set("login", t.Options["login"])
	v.Set("pass", t.Options["password"])
	v.Set("to", phone)
	v.Set("from", t.Options["from"])
	v.Set("text", message)

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

	return nil
}
