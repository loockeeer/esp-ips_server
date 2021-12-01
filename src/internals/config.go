package internals

import (
	"encoding/json"
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

var RssiBufferSize = utils.Atoi("RSSI_BUFFER_SIZE", "RSSI Buffer size is not a number !")

var CONFIG_FILE = os.Getenv("CONFIG_FILE")

var devicesCache []Device

func ListDevices() (devices []Device, err error) {
	if devicesCache != nil {
		return devicesCache, nil
	}

	if CONFIG_FILE == "" {
		CONFIG_FILE = "config.json"
	}

	file, err := os.ReadFile(CONFIG_FILE)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &devices)
	if err != nil {
		return nil, err
	}
	devicesCache = devices
	return devices, nil
}
