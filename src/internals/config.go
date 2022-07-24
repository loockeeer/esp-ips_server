package internals

import (
	"encoding/json"
	"errors"
	"espips_server/src/utils"
	"log"
	"os"
)

// MQTT related config
var (
	MqttHost = os.Getenv("MQTT_HOST")
	MqttPort = utils.Atoi(os.Getenv("MQTT_PORT"), "MQTT port should be number !")
)

// Graphql related config
var (
	ApiHost = os.Getenv("API_HOST")
	ApiPort = utils.Atoi(os.Getenv("API_PORT"), "API port should be number !")
)

// Influx related config
var (
	InfluxHost   = os.Getenv("INFLUX_HOST")
	InfluxPort   = utils.Atoi(os.Getenv("INFLUX_PORT"), "Influx port should be number !")
	InfluxToken  = os.Getenv("INFLUX_TOKEN")
	InfluxOrg    = os.Getenv("INFLUX_ORG")
	InfluxBucket = os.Getenv("INFLUX_BUCKET")
)

var RssiBufferSize = utils.Atoi(os.Getenv("RSSI_BUFFER_SIZE"), "RSSI Buffer size should be number !")
var InitRssiBufferSize = utils.Atoi(os.Getenv("INIT_RSSI_BUFFER_SIZE"), "Init RSSI Buffer size should be number !")

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
	for i, device := range devices {
		if device.Address == nil {
			return nil, errors.New("Missing address for device at index " + string(rune(i)))
		}
		if device.Type == nil {
			return nil, errors.New("Missing type for device " + *device.Address)
		}
		if *device.Type == StationType {
			if device.X == nil {
				return nil, errors.New("Missing X value for station " + *device.Address)
			}
			if device.Y == nil {
				return nil, errors.New("Missing Y value for station " + *device.Address)
			}
		}
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
		if *device.Address == address {
			return &device
		}
	}
	return nil
}
