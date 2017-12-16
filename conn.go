package libmqtt

import "bytes"

// ConPacket is the first packet sent by Client to Server
type ConPacket struct {
	protoName    string
	protoLevel   ProtocolLevel
	Username     string
	Password     string
	ClientId     string
	CleanSession bool
	IsWill       bool
	WillQos      QosLevel
	WillRetain   bool
	Keepalive    uint16
	WillTopic    string
	WillMessage  []byte
}

func (c *ConPacket) Type() CtrlType {
	return CtrlConn
}

func (c *ConPacket) Bytes(buffer *bytes.Buffer) (err error) {
	if buffer == nil || c == nil {
		return
	}
	// fixed header
	// 0x01 0x00
	buffer.WriteByte(CtrlConn << 4)

	payload := c.payload()
	// remaining length
	encodeRemainLength(10+payload.Len(), buffer)

	// Protocol Name and level
	// 0x00 0x04 'M' 'Q' 'T' 'T' 0x04
	buffer.Write([]byte{0x00, 0x04})
	buffer.Write(mqtt)
	buffer.WriteByte(V311)

	// connect flags
	buffer.WriteByte(c.flags())

	// keepalive
	buffer.WriteByte(byte(c.Keepalive >> 8))
	buffer.WriteByte(byte(c.Keepalive))

	_, err = payload.WriteTo(buffer)

	return
}

func (c *ConPacket) flags() byte {
	var connectFlag byte = 0
	if c.ClientId == "" {
		c.CleanSession = true
	}

	if c.CleanSession {
		connectFlag |= 0x02
	}

	if c.IsWill {
		connectFlag |= 0x04
		connectFlag |= c.WillQos << 3

		if c.WillRetain {
			connectFlag |= 0x20
		}
	}

	if c.Password != "" {
		connectFlag |= 0x40
	}

	if c.Username != "" {
		connectFlag |= 0x80
	}

	return connectFlag
}

func (c *ConPacket) payload() *bytes.Buffer {
	result := &bytes.Buffer{}
	// client id
	encodeDataWithLen([]byte(c.ClientId), result)

	// will topic and message
	if c.IsWill {
		encodeDataWithLen([]byte(c.WillTopic), result)
		encodeDataWithLen([]byte(c.WillMessage), result)
	}

	if c.Username != "" {
		encodeDataWithLen([]byte(c.Username), result)
	}

	if c.Password != "" {
		encodeDataWithLen([]byte(c.Password), result)
	}

	return result
}

// ConAckPacket is the packet sent by the Server in response to a ConPacket
// received from a Client.
//
// The first packet sent from the Server to the Client MUST be a ConAckPacket
type ConAckPacket struct {
	Present bool
	Code    ConAckCode
}

func (c *ConAckPacket) Type() CtrlType {
	return CtrlConnAck
}

func (c *ConAckPacket) Bytes(buffer *bytes.Buffer) (err error) {
	if buffer == nil || c == nil {
		return
	}
	// fixed header
	// 0x02 0x00
	buffer.WriteByte(CtrlConnAck << 4)
	buffer.WriteByte(0x02)
	// present flag
	buffer.WriteByte(boolToByte(c.Present))

	// response code
	return buffer.WriteByte(c.Code)
}

// disConPacket is the final Control Packet sent from the Client to the Server.
// It indicates that the Client is disconnecting cleanly.
var (
	DisConPacket = &disConPacket{}
)

type disConPacket struct {
}

func (s *disConPacket) Type() CtrlType {
	return CtrlDisConn
}

func (s *disConPacket) Bytes(buffer *bytes.Buffer) (err error) {
	if buffer == nil || s == nil {
		return
	}
	// fixed header
	buffer.WriteByte(CtrlDisConn << 4)
	return buffer.WriteByte(0x00)
}
