package mqtt

import (
	"errors"
	"espips_server/src/database"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"strconv"
	"strings"
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
	address, payload, err := getMessageInfo(message)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Received RSSI (%s) from %s\n", payload, address)
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

	database.Connection.PushBattery(address, batteryLevel)
}
