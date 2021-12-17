package mqtt

import (
	"errors"
	"espips_server/src/api"
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
	rssiBuffer map[string]map[string][]int
)

func onConnect(client mqtt.Client) {
	log.Println("MQTT Connected")
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
		}
		if *dev.Type == internals.AntennaType {
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

func rssiHandler(client mqtt.Client, message mqtt.Message) {
	sender, payload, err := getMessageInfo(message)
	if err != nil {
		log.Panicln(err)
	}
	log.Printf("Received RSSI (%s) from %s\n", payload, sender)

	scanned := strings.Split(payload, ",")[0]
	rssi := utils.Atoi(strings.Split(payload, ",")[1], "RSSI value received is not a number !")
	rssiBuffer[sender][scanned] = append(rssiBuffer[sender][scanned], rssi)
	err = database.Connection.PushRSSI(sender, scanned, rssi, "")

	go func() {
		switch internals.AppState {
		case internals.INIT_STATE:
			var trainData map[float64]float64
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

						dist := utils.Distance(scannerDev.GetX(), scannerDev.GetY(), scannedDev.GetX(), scannedDev.GetY())

						for _, rssi := range values {
							trainData[float64(rssi)] = dist
						}
					}
				}
			}
			internals.DistanceRssi, err = internals.DistanceRssi.Optimize(trainData)

			if err != nil {
				log.Panicln(err)
			}
			internals.AppState = internals.RUN_STATE
			api.ChangeAppState <- internals.AppState
			break
		case internals.RUN_STATE:
			var data map[internals.Position]float64
			for sender, entries := range rssiBuffer {
				senderDevice := internals.GetDevice(sender)
				if senderDevice == nil {
					continue
				}
				scannerPos := internals.Position{
					X: senderDevice.GetX(),
					Y: senderDevice.GetY(),
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
					data[scannerPos] = internals.DistanceRssi.Execute(avgRSSI)
					rssiBuffer[sender][scanned] = nil
				}
			}
			pos, err := internals.GetPosition(data)
			if err != nil {
				log.Panicln(err)
			}

			err = database.Connection.PushPosition(scanned, database.Position(*pos))
			if err != nil {
				log.Panicln(err)
			}

			scannedDevice := internals.GetDevice(scanned)
			api.PositionEmitter <- internals.GraphQLDevice{
				Address:      *scannedDevice.Address,
				FriendlyName: *scannedDevice.FriendlyName,
				X:            pos.X,
				Y:            pos.Y,
				Speed:        scannedDevice.GetSpeed(),
				Battery:      scannedDevice.GetBattery(),
				Type:         int(*scannedDevice.Type),
			}

			break
		}

	}()
	if err != nil {
		log.Panicln(err)
	}
}

func batteryHandler(client mqtt.Client, message mqtt.Message) {
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
