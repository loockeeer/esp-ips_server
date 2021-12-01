package main

import (
	"espips_server/src/api"
	"espips_server/src/database"
	"espips_server/src/internals"
	"espips_server/src/mqtt"
	"log"
	"sync"
)

var wg = &sync.WaitGroup{}

func main() {
	log.Println("Loading ESP IPS server")

	devices, err := internals.ListDevices()
	if err != nil {
		log.Panicln(err)
	}
	for _, device := range devices {
		log.Printf("Loaded device %s (AKA %s)\n", device.Address, device.FriendlyName)
	}

	log.Println("Connecting to InfluxDB on")
	database.Connect(
		internals.InfluxHost,
		internals.InfluxPort,
		internals.InfluxToken,
		internals.InfluxOrg,
		internals.InfluxBucket)

	log.Println("Connected to InfluxDB")

	// Connect to MQTT Broker
	log.Println("Connecting to mqtt broker")
	wg.Add(1)
	go func() {
		defer wg.Done()
		mqtt.Connect(internals.MqttHost, internals.MqttPort)
	}()

	// Start GraphQL API
	log.Println("Starting GraphQL API")
	wg.Add(1)
	go func() {
		defer wg.Done()
		api.Start(internals.ApiHost, internals.ApiPort)
	}()

	wg.Wait()
}
