package api

import (
	"espips_server/src/internals"
	"fmt"
	pahoMqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
)

var client pahoMqtt.Client

func ConnectMQTT(host string, port int) {
	opts := pahoMqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", host, port))
	opts.SetClientID("server")
	opts.OnConnect = onConnect

	client = pahoMqtt.NewClient(opts)

	log.Println("Connected to broker")

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Panicln(token)
	}
}

func GlobalControl(state int) {
	devices, _ := internals.ListDevices()
	for _, device := range devices {
		payload := ""

		/*
			#define CARMODE 0
			#define ANTENNAMODE_INIT 1
			#define ANTENNAMODE_RUN 2
			#define IDLE_MODE 3
		*/
		if state == internals.IDLE_STATE {
			payload = "3"
		}
		if *device.Type == internals.AntennaType {
			if state == internals.RUN_STATE {
				payload = "2"
			} else if state == internals.INIT_STATE {
				payload = "1"
			}
		} else if *device.Type == internals.CarType {
			if state == internals.RUN_STATE {
				payload = "0"
			} else {
				payload = "3"
			}
		}
		client.Publish(fmt.Sprintf("cc/%s", *device.Address), 2, false, payload)
	}
}
