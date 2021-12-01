package main

import (
	"espips_server/src/api"
	"espips_server/src/database"
	"espips_server/src/mqtt"
	"log"
	"sync"
)

var wg = &sync.WaitGroup{}

func main() {
	log.Println("Loading ESP IPS server")

	log.Println("Connecting to InfluxDB on")
	database.Connect(
		INFLUX_HOST,
		INFLUX_PORT,
		INFLUX_TOKEN,
		INFLUX_ORG,
		INFLUX_BUCKET)

	log.Println("Connected to InfluxDB")

	// Connect to MQTT Broker
	log.Println("Connecting to mqtt broker")
	wg.Add(1)
	go func() {
		defer wg.Done()
		mqtt.Connect(MQTT_HOST, MQTT_PORT)
	}()

	// Start GraphQL API
	log.Println("Starting GraphQL API")
	wg.Add(1)
	go func() {
		defer wg.Done()
		api.Start(API_HOST, API_PORT)
	}()

	wg.Wait()
}
