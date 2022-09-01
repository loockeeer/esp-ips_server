package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/loockeeer/espipsserver/src/communication"
	"github.com/loockeeer/espipsserver/src/ips"
	"github.com/loockeeer/espipsserver/src/tools"
	"log"
	"os"
	"strconv"
	"time"
)

func appendCSVRow(fileName string, row []string) error {
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	w := csv.NewWriter(f)
	err = w.Write(row)
	if err != nil {
		return err
	}
	w.Flush()
	return nil
}

type Device struct {
	FriendlyName string   `json:"friendlyName"`
	Address      string   `json:"address"`
	Type         string   `json:"type"`
	X            *float64 `json:"x,omitempty"`
	Y            *float64 `json:"y,omitempty"`
}

type Config struct {
	Model   *ips.DistanceRssiModel `json:"model,omitempty"`
	Devices []Device               `json:"devices"`
}

var ConfigFile string
var Mode string
var RssiMinTrain int
var RssiMaxRun int
var MqttHost string
var MqttPort int
var MaxTimeDifference time.Duration
var IPSConfig Config

func init() {
	var err error
	ConfigFile = os.Getenv("CONFIG_FILE")
	Mode = os.Getenv("MODE")
	RssiMinTrain, err = strconv.Atoi(os.Getenv("RSSI_MIN_TRAIN"))
	if err != nil {
		log.Fatalln(err)
	}

	RssiMaxRun, err = strconv.Atoi(os.Getenv("RSSI_MAX_RUN"))
	if err != nil {
		log.Fatalln(err)
	}
	MqttHost = os.Getenv("MQTT_HOST")
	MqttPort, err = strconv.Atoi(os.Getenv("MQTT_PORT"))
	if err != nil {
		log.Fatalln(err)
	}

	MaxTimeDifference, err = time.ParseDuration(os.Getenv("MAX_TIME_DIFFERENCE"))
	if err != nil {
		log.Fatalln(err)
	}

	if (Mode != "run") && (Mode != "train") && (Mode != "debug") {
		log.Fatalln("MODE should be one of run,train,debug")
	}

	file, err := os.ReadFile(ConfigFile)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(file, &IPSConfig)
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	api := communication.NewDeviceAPI(communication.DeviceAPIOptions{
		MqttHost: MqttHost,
		MqttPort: MqttPort,
	})
	go func() {
		token := api.Start()
		log.Println("Started device API")
		if token.Wait() && token.Error() != nil {
			log.Panicln(token)
		}
	}()
	if Mode == "train" {
		trainMain(api)
	} else if Mode == "run" {
		runMain(api)
	} else if Mode == "debug" {
		debugMain(api)
	}
}

func debugMain(api communication.DeviceAPI) {
	distances := getDistances(IPSConfig.Devices)
	api.OnRSSI(func(event *communication.RSSIEvent) {
		if tools.Find[Device](IPSConfig.Devices, func(d Device) bool {
			return d.Address == event.Scanned.Address
		}) == nil {
			return
		}
		// DEBUG CODE
		err := appendCSVRow("debug.csv", []string{event.Scanner.Address, event.Scanned.Address, strconv.Itoa(event.RSSI), fmt.Sprintf("%f", distances[event.Scanner.Address][event.Scanned.Address])}) // SCANNER,SCANNED,RSSI,DISTANCE
		if err != nil {
			log.Fatalln(err)
		}
	})

}

func runMain(api communication.DeviceAPI) {
	positions := map[string]ips.Position{}
	for _, device := range IPSConfig.Devices {
		if device.Type == "station" {
			positions[device.Address] = ips.Position{
				X: *device.X,
				Y: *device.Y,
			}
			api.Send(device.Address, communication.DeviceStationRunMode)
		} else if device.Type == "beacon" {
			api.Send(device.Address, communication.DeviceBeaconMode)
		}
	}

	wrapper := ips.NewIPSWrapper(RssiMaxRun, RssiMaxRun, MaxTimeDifference, *IPSConfig.Model, positions)
	wrapper.OnPosition(func(event *ips.PositionEvent) {
		log.Printf("Got position for %s [X = %f][Y = %f]", event.Device, event.Position.X, event.Position.Y)
	})
	api.OnRSSI(func(event *communication.RSSIEvent) {
		if tools.Find[Device](IPSConfig.Devices, func(d Device) bool {
			return d.Address == event.Scanned.Address
		}) == nil {
			return
		}
		err := wrapper.Collect(event.Scanner.Address, event.Scanned.Address, event.RSSI)
		if err != nil {
			log.Fatalln(err)
		}
		// DEBUG CODE
		err = appendCSVRow("debug.csv", []string{event.Scanner.Address, event.Scanned.Address, strconv.Itoa(event.RSSI), fmt.Sprintf("%f", IPSConfig.Model.Execute(float64(event.RSSI)))}) // SCANNER,SCANNED,RSSI,DISTANCE
		if err != nil {
			log.Fatalln(err)
		}
	})
}

func getDistances(config []Device) map[string]map[string]float64 {
	var distances = map[string]map[string]float64{}
	for _, d1 := range config {
		for _, d2 := range config {
			dist := ips.Distance(ips.Position{X: *d1.X, Y: *d1.Y}, ips.Position{X: *d2.X, Y: *d2.Y})
			if _, ok := distances[d1.Address]; ok {
				distances[d1.Address] = map[string]float64{
					d2.Address: dist,
				}
			} else {
				distances[d1.Address][d2.Address] = dist
			}
		}
	}
	return distances
}

func trainMain(api communication.DeviceAPI) {
	for _, device := range IPSConfig.Devices {
		if device.Type == "station" {
			api.Send(device.Address, communication.DeviceStationInitMode)
		} else if device.Type == "beacon" {
			api.Send(device.Address, communication.DeviceIdleMode)
		}
	}
	rssiCollector := ips.NewRSSICollector(RssiMinTrain, time.Hour*24*365)
	api.OnRSSI(func(event *communication.RSSIEvent) {
		if tools.Find[Device](IPSConfig.Devices, func(d Device) bool {
			return d.Address == event.Scanned.Address
		}) == nil {
			return
		}
		rssiCollector.Collect(event.Scanner.Address, event.Scanned.Address, event.RSSI)
		for _, configDevice := range IPSConfig.Devices {
			if _, ok := rssiCollector.Data[configDevice.Address]; !ok {
				return
			} else {
				for _, configDevice2 := range IPSConfig.Devices {
					if queue, ok := rssiCollector.Data[configDevice.Address][configDevice2.Address]; !ok {
						return
					} else {
						if len(queue.Data) < RssiMinTrain {
							return
						}
					}
				}
			}
		}
		model := ips.NewDistanceRssiModel()
		err := model.Train(rssiCollector, getDistances(IPSConfig.Devices))
		if err != nil {
			log.Fatalln(err)
		}

		IPSConfig.Model = &model
		data, err := json.Marshal(IPSConfig)
		if err != nil {
			log.Fatalln(err)
		}

		err = os.WriteFile(ConfigFile, data, 0644)
		if err != nil {
			log.Fatalln(err)
		}
	})
}
