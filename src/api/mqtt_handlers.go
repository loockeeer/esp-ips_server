package api

import (
	"errors"
	"espips_server/src/database"
	"espips_server/src/internals"
	"espips_server/src/utils"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"strconv"
	"strings"
)

var (
	rssiBuffer = map[string]map[string][]int{}
)

func onConnect(client mqtt.Client) {
	log.Println("MQTT Connected")
	client.Subscribe("cc/+", 2, ccHandler)
	client.Subscribe("rssi/+", 2, rssiHandler)
	client.Subscribe("battery/+", 2, batteryHandler)
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
		log.Panicln(err)
	}
	switch payload {
	case "4":
		log.Printf("Received ack from %s\n", address)
		dev := internals.GetDevice(address)
		if dev == nil {
			break
		}

		payload := ""
		if internals.AppState == internals.IDLE_STATE {
			payload = "3"
		} else if *dev.Type == internals.AntennaType {
			if internals.AppState == internals.RUN_STATE {
				payload = "2"
			} else if internals.AppState == internals.INIT_STATE {
				payload = "1"
			}
		} else if *dev.Type == internals.CarType {
			if internals.AppState == internals.RUN_STATE {
				payload = "0"
			} else {
				payload = "3"
			}
		}

		client.Publish(fmt.Sprintf("cc/%s", address), 2, false, payload)

		break
	}
}

func rssiHandler(_ mqtt.Client, message mqtt.Message) {
	sender, payload, err := getMessageInfo(message)
	if err != nil {
		log.Panicln(err)
	}

	scanned := strings.Split(payload, ",")[0]
	rssi, err := strconv.Atoi(strings.Split(payload, ",")[1])
	if err != nil {
		log.Println(err)
		return
	}
	if rssiBuffer[sender] == nil {
		rssiBuffer[sender] = map[string][]int{}
	}
	rssiBuffer[sender][scanned] = append(rssiBuffer[sender][scanned], rssi)
	_ = database.Connection.PushRSSI(sender, scanned, rssi, "")
	switch internals.AppState {
	case internals.INIT_STATE:
		var trainData = map[float64]float64{}
		for scanner, entries := range rssiBuffer {
			for scanned, values := range entries {
				if len(values) < internals.InitRssiBufferSize {
					return
				} else {
					scannerDev := internals.GetDevice(scanner)
					scannedDev := internals.GetDevice(scanned)

					if *scannerDev.Type != internals.AntennaType {
						break
					} else if *scannedDev.Type != internals.AntennaType {
						continue
					}
					dist := utils.Distance(scannerDev.GetPosition().X, scannerDev.GetPosition().Y, scannedDev.GetPosition().X, scannedDev.GetPosition().Y)

					for _, rssi := range values {
						trainData[float64(rssi)] = dist
					}
				}
			}
		}
		internals.DistanceRssi, err = internals.DistanceRssi.Optimize(trainData)

		if err != nil {
			log.Println("Failed to init distance-rssi model : ", err)
		} else {
			log.Println(internals.DistanceRssi)
		}
		internals.AppState = internals.RUN_STATE
		GlobalControl(internals.AppState)
		rssiBuffer = map[string]map[string][]int{}
		break
	case internals.RUN_STATE:
		var data = map[internals.Position]float64{}
		scannedDevice := internals.GetDevice(scanned)
		if *scannedDevice.Type != 1 {
			break
		}
		for sender, entries := range rssiBuffer {
			senderDevice := internals.GetDevice(sender)
			if senderDevice == nil {
				continue
			}
			for _scanned, values := range entries {
				if _scanned != scanned {
					continue
				}
				if len(values) < internals.RssiBufferSize {
					return
				}
				avgRSSI := 0.0
				for _, rssi := range values {
					avgRSSI += float64(rssi)
				}
				avgRSSI /= float64(len(values))
				data[*senderDevice.GetPosition()] = internals.DistanceRssi.Execute(avgRSSI)
				log.Println(avgRSSI, data[*senderDevice.GetPosition()])
				rssiBuffer[sender][scanned] = nil
			}
		}
		pos, err := internals.GetPosition(data)
		if err != nil {
			log.Panicln(err)
		}

		err = scannedDevice.SetPosition(pos)
		if err != nil {
			log.Panicln(err)
		}
		PositionEvent.Emit(internals.GraphQLDevice{
			Address:      *scannedDevice.Address,
			FriendlyName: *scannedDevice.FriendlyName,
			X:            pos.X,
			Y:            pos.Y,
			Speed:        scannedDevice.GetSpeed(),
			Battery:      scannedDevice.GetBattery(),
			Type:         int(*scannedDevice.Type),
		})
		log.Printf("X : %f | Y : %f\n", pos.X, pos.Y)

		break
	}
	if err != nil {
		log.Panicln(err)
	}
}

func batteryHandler(_ mqtt.Client, message mqtt.Message) {
	address, payload, err := getMessageInfo(message)
	if err != nil {
		log.Panicln(err)
	}
	log.Printf("Received battery info (%s) from %s\n", payload, address)

	batteryLevel, err := strconv.ParseFloat(payload, 64)
	if err != nil {
		log.Panicln(err)
	}

	err = database.Connection.PushBattery(address, batteryLevel)
	if err != nil {
		log.Panicln(err)
	}
}
