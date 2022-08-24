package main

import (
	"github.com/loockeeer/espipsserver/src/communication"
	"github.com/loockeeer/espipsserver/src/ips"
	"log"
)

func main() {
	api := communication.NewDeviceAPI(communication.DeviceAPIOptions{
		MqttHost: "",
		MqttPort: 1111,
	})
	go func() {
		token := api.Start()
		log.Println("Started device API")
		if token.Wait() && token.Error() != nil {
			log.Panicln(token)
		}
	}()

}

func trainMain(api communication.DeviceAPI) {
	rssiCollector := ips.NewRSSICollector(3)
	api.OnRSSI(func(event *communication.RSSIEvent) {
		rssiCollector.Collect(event.Scanner.Address, event.Scanned.Address, event.RSSI)
	})
}
