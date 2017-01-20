package main

import (
	"github.com/denniselite/sms-service/manager"
	"github.com/denniselite/sms-service/structs"
	"github.com/denniselite/go-toolkit/api"
	"github.com/denniselite/go-toolkit/conn"
	. "github.com/denniselite/go-toolkit/errors"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile | log.LUTC)

	config := loadConfig()

	// get rabbit connection
	var rmq *conn.Rmq

	rabbitConnectionString := fmt.Sprintf("amqp://%s:%s@%s:%d",
		config.Rabbit.Username,
		config.Rabbit.Password,
		config.Rabbit.Host,
		config.Rabbit.Port,
	)

	man := new(manager.SmsManager)

	// Load sms provider settings
	man.SmsProviders = config.Sms

	run := func() {
		man.Run(rmq)
	}

	rmq, err := conn.NewRmq(rabbitConnectionString, run)
	OopsT(api.SysToken, err)

	run()

	// setup api server
	log.Printf("%s Listen HTTP port: %d\n", api.SysToken, config.Listen)
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "pong")
	})
	OopsT(api.SysToken, http.ListenAndServe(fmt.Sprintf(":%d", config.Listen), nil))
}

func loadConfig() structs.Config {
	var filename string

	// register flags
	flag.StringVar(&filename, "config", "", "config filename")
	flag.StringVar(&filename, "c", "", "config filename (shorthand)")

	flag.Parse()

	config := structs.Config{}

	configData, err := ioutil.ReadFile(filename)
	OopsT(api.SysToken, err)
	err = yaml.Unmarshal(configData, &config)
	OopsT(api.SysToken, err)

	return config
}
