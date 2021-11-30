package mqtt

import (
	"fmt"
	pahoMqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
)

func Connect(host string, port string) {
	opts := pahoMqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", host, port))
	opts.SetClientID("server")
	opts.OnConnect = onConnect

	client := pahoMqtt.NewClient(opts)

	client.Subscribe("cc/*", 2, ccHandler)
	client.Subscribe("rssi/*", 2, rssiHandler)
	client.Subscribe("battery/*", 2, batteryHandler)

	log.Println("Connected to broker")

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}
