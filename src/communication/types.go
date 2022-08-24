package communication

type DeviceMode = int

const (
	DeviceIdleMode DeviceMode = iota
	DeviceStationInitMode
	DeviceStationRunMode
	DeviceBeaconMode
)

type MQTTFormat struct {
	Topic   string
	Payload string
	QoS     byte
}

type DeviceAPIOptions struct {
	MqttHost     string
	MqttPort     int
	MqttPassword string
	MqttUsername string
	MqttId       string
	PingInterval int
	KeepAlive    int
}

type ConnectEvent struct {
	Device Device
}
type DisconnectEvent struct {
	Device Device
}
type RSSIEvent struct {
	Scanner Device
	Scanned Device
	RSSI    int
}
type PingEvent struct {
	Device Device
}
