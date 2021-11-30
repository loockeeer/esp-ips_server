package main

import (
	"espips_server/src/api"
	"espips_server/src/database"
	"espips_server/src/mqtt"
	"log"
	"os"
	"sync"
)

var wg = &sync.WaitGroup{}

func main() {
	log.Println("Loading ESP IPS server")

	log.Println("Connecting to InfluxDB on")
	database.Connect(
		os.Getenv("INFLUX_HOST"),
		os.Getenv("INFLUX_PORT"),
		os.Getenv("INFLUX_TOKEN"),
		os.Getenv("INFLUX_ORG"),
		os.Getenv("INFLUX_BUCKET"))

	log.Println("Connected to InfluxDB")

	// Connect to MQTT Broker
	log.Println("Connecting to mqtt broker")
	wg.Add(1)
	go func() {
		defer wg.Done()
		mqtt.Connect(os.Getenv("MQTT_HOST"), os.Getenv("MQTT_PORT"))
	}()

	// Start GraphQL API
	log.Println("Starting GraphQL API")
	wg.Add(1)
	go func() {
		defer wg.Done()
		api.Start(os.Getenv("API_HOST"), os.Getenv("API_PORT"))
	}()

	wg.Wait()
}
