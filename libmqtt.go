package libmqtt

type QosLevel = byte

type ProtocolLevel = byte

const (
	V31 ProtocolLevel = iota + 3
	V311
)

const (
	Qos0 QosLevel = iota
	Qos1
	Qos2
)

var (
	MQTT = []byte{'M', 'Q', 'T', 'T'}
)
