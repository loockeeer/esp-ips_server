package main

import (
	"espips_server/src/internals"
	"espips_server/src/utils"
	"os"
)

var (
	MQTT_HOST = os.Getenv("MQTT_HOST")
	MQTT_PORT = utils.Atoi("MQTT_PORT", "MQTT port is not a number !")
)

var (
	API_HOST = os.Getenv("API_HOST")
	API_PORT = utils.Atoi("API_PORT", "API port is not a number !")
)

var (
	INFLUX_HOST   = os.Getenv("INFLUX_HOST")
	INFLUX_PORT   = utils.Atoi("INFLUX_PORT", "Influx port is not a number !")
	INFLUX_TOKEN  = os.Getenv("INFLUX_TOKEN")
	INFLUX_ORG    = os.Getenv("INFLUX_ORG")
	INFLUX_BUCKET = os.Getenv("INFLUX_BUCKET")
)

func ListDevices() []internals.Device {
	return nil
}
