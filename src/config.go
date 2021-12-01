package main

import (
	"espips_server/src/internals"
	"espips_server/src/utils"
	"os"
)

var (
	MqttHost = os.Getenv("MQTT_HOST")
	MqttPort = utils.Atoi("MQTT_PORT", "MQTT port is not a number !")
)

var (
	ApiHost = os.Getenv("API_HOST")
	ApiPort = utils.Atoi("API_PORT", "API port is not a number !")
)

var (
	InfluxHost   = os.Getenv("INFLUX_HOST")
	InfluxPort   = utils.Atoi("INFLUX_PORT", "Influx port is not a number !")
	InfluxToken  = os.Getenv("INFLUX_TOKEN")
	InfluxOrg    = os.Getenv("INFLUX_ORG")
	InfluxBucket = os.Getenv("INFLUX_BUCKET")
)

func ListDevices() []internals.Device {
	return nil
}
