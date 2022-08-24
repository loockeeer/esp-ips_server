package communication

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/loockeeer/espipsserver/src/tools"
	"strconv"
)

type Device struct {
	Address string
	wrapper DeviceAPI
}

func (d *Device) Send(mode DeviceMode) {
	d.wrapper.Send(d.Address, mode)
}

func GetPayloadForMode(address string, mode DeviceMode) MQTTFormat {
	return MQTTFormat{
		Topic:   "cc/" + address,
		Payload: strconv.Itoa(mode),
		QoS:     2,
	}
}

func GetPayloadForBroadcast(mode DeviceMode) MQTTFormat {
	return MQTTFormat{
		Topic:   "cc",
		Payload: strconv.Itoa(mode),
		QoS:     2,
	}
}

type DeviceAPI struct {
	MqttHost          string
	MqttPort          int
	MqttPassword      string
	MqttUsername      string
	MqttId            string
	mqttClient        mqtt.Client
	PingInterval      int
	KeepAlive         int
	devices           []Device
	connectHandler    tools.EventHandler[ConnectEvent]
	disconnectHandler tools.EventHandler[DisconnectEvent]
	rssiHandler       tools.EventHandler[RSSIEvent]
	pingHandler       tools.EventHandler[PingEvent]
}

func NewDeviceAPI(options DeviceAPIOptions) DeviceAPI {
	return DeviceAPI{
		MqttHost:          options.MqttHost,
		MqttPort:          options.MqttPort,
		MqttPassword:      options.MqttPassword,
		MqttUsername:      options.MqttUsername,
		MqttId:            options.MqttId,
		PingInterval:      options.PingInterval,
		KeepAlive:         options.KeepAlive,
		mqttClient:        nil,
		devices:           []Device{},
		connectHandler:    tools.NewEventHandler[ConnectEvent](),
		disconnectHandler: tools.NewEventHandler[DisconnectEvent](),
		rssiHandler:       tools.NewEventHandler[RSSIEvent](),
		pingHandler:       tools.NewEventHandler[PingEvent](),
	}
}

func (w *DeviceAPI) NewDevice(address string) Device {
	return Device{
		Address: address,
		wrapper: *w,
	}
}

func (w *DeviceAPI) Send(address string, mode DeviceMode) mqtt.Token {
	format := GetPayloadForMode(address, mode)
	return w.mqttClient.Publish(format.Topic, format.QoS, false, format.Payload)
}

func (w *DeviceAPI) BroadcastAll(mode DeviceMode) mqtt.Token {
	format := GetPayloadForBroadcast(mode)
	return w.mqttClient.Publish(format.Topic, format.QoS, false, format.Payload)
}

func (w *DeviceAPI) OnConnect(listener func(*ConnectEvent)) {
	w.connectHandler.Register(listener)
}

func (w *DeviceAPI) OnDisconnect(listener func(*DisconnectEvent)) {
	w.disconnectHandler.Register(listener)
}

func (w *DeviceAPI) OnRSSI(listener func(*RSSIEvent)) {
	w.rssiHandler.Register(listener)
}

func (w *DeviceAPI) OnPing(listener func(*PingEvent)) {
	w.pingHandler.Register(listener)
}

func (w *DeviceAPI) Start() mqtt.Token {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", w.MqttHost, w.MqttPort))
	if w.MqttId != "" {
		opts.SetClientID(w.MqttId)
	} else {
		opts.SetClientID("espips_server")
	}

	if w.MqttUsername != "" && w.MqttPassword != "" {
		opts.SetUsername(w.MqttUsername)
		opts.SetPassword(w.MqttPassword)
	}
	format := GetPayloadForBroadcast(DeviceIdleMode)
	opts.SetWill(format.Topic, format.Payload, format.QoS, false)
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		client.Subscribe("rssi/+", 2, func(client mqtt.Client, message mqtt.Message) {
			go RSSIHandler(w, client, message)
		})
		client.Subscribe("ping/+", 2, func(client mqtt.Client, message mqtt.Message) {
			go PingHandler(w, client, message)
		})
	})
	w.mqttClient = mqtt.NewClient(opts)

	return w.mqttClient.Connect()
}
