package internals

import (
	"encoding/json"
	"espips_server/src/utils"
	"log"
	"os"
)

var (
	MqttHost = os.Getenv("MQTT_HOST")
	MqttPort = utils.Atoi(os.Getenv("MQTT_PORT"), "MQTT port should be number !")
)

var (
	ApiHost = os.Getenv("API_HOST")
	ApiPort = utils.Atoi(os.Getenv("API_PORT"), "API port should be number !")
)

var (
	InfluxHost   = os.Getenv("INFLUX_HOST")
	InfluxPort   = utils.Atoi(os.Getenv("INFLUX_PORT"), "Influx port should be number !")
	InfluxToken  = os.Getenv("INFLUX_TOKEN")
	InfluxOrg    = os.Getenv("INFLUX_ORG")
	InfluxBucket = os.Getenv("INFLUX_BUCKET")
)

var RssiBufferSize = utils.Atoi(os.Getenv("RSSI_BUFFER_SIZE"), "RSSI Buffer size should be number !")
var InitRssiBufferSize = utils.Atoi(os.Getenv("INIT_RSSI_BUFFER_SIZE"), "Init RSSI Buffer size should be number !")

var RssiDistanceOrder = utils.Atoi(os.Getenv("RSSI_DISTANCE_ORDER"), "RSSI-Distance relation order should be number !")

var ConfigFile = os.Getenv("CONFIG_FILE")

var devicesCache []Device

func ListDevices() (devices []Device, err error) {
	if devicesCache != nil {
		return devicesCache, nil
	}

	if ConfigFile == "" {
		ConfigFile = "config.json"
	}

	file, err := os.ReadFile(ConfigFile)
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

func GetDevice(address string) *Device {
	devices, err := ListDevices()
	if err != nil {
		log.Panicln(err)
	}
	for _, device := range devices {
		if device.Address == address {
			return &device
		}
	}
	return &Device{}
}
