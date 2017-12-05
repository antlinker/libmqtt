package libmqtt

import "bytes"

type ConnPacket struct {
	username string
	password string
	proto    ProtocolLevel

	ClientId     string
	CleanSession bool
	IsWill       bool
	WillQos      QosLevel
	WillRetain   bool
	Keepalive    uint16

	WillTopic   string
	WillMessage string
}

func (c *ConnPacket) Type() CtrlType {
	return CtrlConn
}

func (c *ConnPacket) Bytes(buffer *bytes.Buffer) (err error) {
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
	buffer.Write(MQTT)
	buffer.WriteByte(c.proto)

	// connect flags
	buffer.WriteByte(c.connectFlags())

	// keepalive
	buffer.WriteByte(byte(c.Keepalive >> 8))
	buffer.WriteByte(byte(c.Keepalive))

	_, err = payload.WriteTo(buffer)

	return
}

func (c *ConnPacket) connectFlags() byte {
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

	if c.password != "" {
		connectFlag |= 0x40
	}

	if c.username != "" {
		connectFlag |= 0x80
	}

	return connectFlag
}

func (c *ConnPacket) payload() *bytes.Buffer {
	result := &bytes.Buffer{}
	// client id
	lenClientId := len(c.ClientId)
	result.WriteByte(byte(lenClientId >> 8))
	result.WriteByte(byte(lenClientId))
	result.Write([]byte(c.ClientId))

	// will topic and message
	if c.IsWill {
		lenWillTopic := len(c.WillTopic)
		result.WriteByte(byte(lenWillTopic >> 8))
		result.WriteByte(byte(lenWillTopic))
		result.Write([]byte(c.WillTopic))

		lenWillMsg := len(c.WillMessage)
		result.WriteByte(byte(lenWillMsg >> 8))
		result.WriteByte(byte(lenWillMsg))
		result.Write([]byte(c.WillMessage))
	}

	if c.username != "" {
		lenUsername := len(c.username)
		result.WriteByte(byte(lenUsername >> 8))
		result.WriteByte(byte(lenUsername))
		result.Write([]byte(c.username))
	}

	if c.password != "" {
		lenPassword := len(c.password)
		result.WriteByte(byte(lenPassword >> 8))
		result.WriteByte(byte(lenPassword))
		result.Write([]byte(c.password))
	}

	return result
}

type ConnAckCode = byte

const (
	ConnAccepted ConnAckCode = iota
	ConnBadProtocol
	ConnIdRejected
	ConnServerUnavailable
	ConnBadIdentity
	ConnAuthFail
)

type ConnAckPacket struct {
	Present bool
	Code    ConnAckCode
}

func (c *ConnAckPacket) Type() CtrlType {
	return CtrlConnAck
}

func (c *ConnAckPacket) Bytes(buffer *bytes.Buffer) (err error) {
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

// DisConnPacket is the final Control Packet sent from the Client to the Server.
// It indicates that the Client is disconnecting cleanly.
type DisConnPacket struct {
}

func (s *DisConnPacket) Type() CtrlType {
	return CtrlDisConn
}

func (s *DisConnPacket) Bytes(buffer *bytes.Buffer) (err error) {
	if buffer == nil || s == nil {
		return
	}
	// fixed header
	buffer.WriteByte(CtrlDisConn << 4)
	return buffer.WriteByte(0x00)
}
