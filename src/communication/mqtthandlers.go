package communication

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"strconv"
	"strings"
)

func RSSIHandler(devicewrapper *DeviceCommunicationWrapper, client mqtt.Client, message mqtt.Message) {
	scannerA := strings.Split(message.Topic(), "/")
	if len(scannerA) != 2 {
		return
	}
	scanner := scannerA[1]
	rssiA := strings.Split(string(message.Payload()), ",")
	if len(rssiA) != 2 {
		return
	}
	scanned := rssiA[0]
	rssi, err := strconv.Atoi(rssiA[1])
	if err != nil {
		panic(err)
	}
	devicewrapper.rssiHandler.Dispatch(&RSSIEvent{
		RSSI:    rssi,
		Scanner: devicewrapper.NewDevice(scanner),
		Scanned: devicewrapper.NewDevice(scanned),
	})
}

func PingHandler(devicewrapper *DeviceCommunicationWrapper, client mqtt.Client, message mqtt.Message) {

}
