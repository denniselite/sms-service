package components

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type ISMSIndonesia struct {
	Options map[string]string
}

func (t *ISMSIndonesia) SetOptions(options map[string]string) {
	t.Options = options
}

func (t *ISMSIndonesia) GetName() (name string) {
	return "iSMSIndonesia"
}

func (t *ISMSIndonesia) SendSms(phone string, message string) error {
	req, err := http.NewRequest("GET", t.Options["url"], nil)
	if err != nil {
		return err
	}

	v := req.URL.Query()
	v.Add("un", t.Options["login"])
	v.Add("pwd", t.Options["password"])
	v.Add("dstno", phone)
	v.Add("msg", message)
	v.Add("type", t.Options["serviceType"])
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

	response := string(contents[:len(contents)])
	if !strings.Contains(response, "SUCCESS") {
		return errors.New(fmt.Sprintf("Error: %s", contents))
	}

	return nil
}
