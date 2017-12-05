package libmqtt

import "bytes"

type Packet interface {
	Type() CtrlType

	// Bytes dump a MQTT Packet object to MQTT bytes into provided buffer
	Bytes(*bytes.Buffer) error
}

type Topic struct {
	Name string
	Qos  QosLevel
}
