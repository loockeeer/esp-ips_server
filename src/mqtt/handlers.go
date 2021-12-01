package mqtt

import (
	"errors"
	"espips_server/src/database"
	"espips_server/src/utils"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"strconv"
	"strings"
)

var (
	rssiBuffer map[string]map[string][]int
)

func onConnect(client mqtt.Client) {
	log.Println("MQTT Connected")
}

func getMessageInfo(message mqtt.Message) (address string, payload string, err error) {
	components := strings.Split(message.Topic(), "/")
	if len(components) < 2 {
		return "", string(message.Payload()), errors.New("Missing address")
	} else {
		return components[1], string(message.Payload()), nil
	}
}

func ccHandler(client mqtt.Client, message mqtt.Message) {
	address, payload, err := getMessageInfo(message)
	if err != nil {
		log.Panic(err)
	}
	switch payload {
	case "4":
		log.Printf("Received ack from %s\n", address)
	case "5":
		log.Printf("Received keep from %s\n", address)
		err := database.Connection.PushKeep(address)
		if err != nil {
			log.Panic(err)
		}
	}
}

func rssiHandler(client mqtt.Client, message mqtt.Message) {
	sender, payload, err := getMessageInfo(message)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Received RSSI (%s) from %s\n", payload, sender)

	scanned := strings.Split(payload, ",")[0]
	rssi := utils.Atoi(strings.Split(payload, ",")[1], "RSSI value received is not a number !")
	rssiBuffer[scanned][sender] = append(rssiBuffer[scanned][sender], rssi)
}

func batteryHandler(client mqtt.Client, message mqtt.Message) {
	address, payload, err := getMessageInfo(message)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Received battery info (%s) from %s\n", payload, address)

	batteryLevel, err := strconv.ParseFloat(payload, 64)
	if err != nil {
		log.Panic(err)
	}

	err = database.Connection.PushBattery(address, batteryLevel)
	if err != nil {
		log.Panicln(err)
	}
}
