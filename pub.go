package libmqtt

import "bytes"

// PublishPacket is sent from a Client to a Server or from Server to a Client
// to transport an Application Message.
type PublishPacket struct {
	IsDup     bool
	Qos       QosLevel
	IsRetain  bool
	TopicName string
	Payload   []byte
	packetId  uint16
}

func (p *PublishPacket) Type() CtrlType {
	return CtrlPublish
}

func (p *PublishPacket) Bytes(buffer *bytes.Buffer) (err error) {
	if buffer == nil || p == nil {
		return
	}
	// fixed header
	buffer.WriteByte(CtrlSubscribe<<4 | boolToByte(p.IsDup)<<3 |
		boolToByte(p.IsRetain) | p.Qos<<1)
	payload := p.payload()
	encodeRemainLength(payload.Len(), buffer)
	_, err = payload.WriteTo(buffer)
	return
}

func (p *PublishPacket) payload() (result *bytes.Buffer) {
	result = &bytes.Buffer{}
	encodeDataWithLen([]byte(p.TopicName), result)
	if p.Qos > Qos0 {
		result.WriteByte(byte(p.packetId >> 8))
		result.WriteByte(byte(p.packetId))
	}
	result.Write(p.Payload)
	return
}

// PubAckPacket is the response to a PublishPacket with QoS level 1.
type PubAckPacket struct {
	PacketId uint16
}

func (p *PubAckPacket) Type() CtrlType {
	return CtrlPubAck
}

func (p *PubAckPacket) Bytes(buffer *bytes.Buffer) (err error) {
	if buffer == nil || p == nil {
		return
	}

	// fixed header
	buffer.WriteByte(CtrlPubComp << 4)
	// remaining length
	buffer.WriteByte(0x02)
	// packet id
	buffer.WriteByte(byte(p.PacketId >> 8))
	return buffer.WriteByte(byte(p.PacketId))
}

// PubRecvPacket is the response to a PublishPacket with QoS 2.
// It is the second packet of the QoS 2 protocol exchange.
type PubRecvPacket struct {
	PacketId uint16
}

func (p *PubRecvPacket) Type() CtrlType {
	return CtrlPubRecv
}

func (p *PubRecvPacket) Bytes(buffer *bytes.Buffer) (err error) {
	if buffer == nil || p == nil {
		return
	}

	// fixed header
	buffer.WriteByte(CtrlPubRecv << 4)
	// remaining length
	buffer.WriteByte(0x02)
	// packet id
	buffer.WriteByte(byte(p.PacketId >> 8))
	return buffer.WriteByte(byte(p.PacketId))
}

// PubRelPacket is the response to a PubRecvPacket.
// It is the third packet of the QoS 2 protocol exchange.
type PubRelPacket struct {
	PacketId uint16
}

func (p *PubRelPacket) Type() CtrlType {
	return CtrlPubRel
}

func (p *PubRelPacket) Bytes(buffer *bytes.Buffer) (err error) {
	if buffer == nil || p == nil {
		return
	}

	buffer.WriteByte(CtrlPubRel<<4 | 0x02)
	// remaining length
	buffer.WriteByte(0x02)
	// packet id
	buffer.WriteByte(byte(p.PacketId >> 8))
	return buffer.WriteByte(byte(p.PacketId))
}

// PubCompPacket is the response to a PubRelPacket.
// It is the fourth and final packet of the QoS 892 2 protocol exchange. 893
type PubCompPacket struct {
	PacketId uint16
}

func (p *PubCompPacket) Type() CtrlType {
	return CtrlPubComp
}

func (p *PubCompPacket) Bytes(buffer *bytes.Buffer) (err error) {
	if buffer == nil || p == nil {
		return
	}
	// fixed header
	buffer.WriteByte(CtrlPubComp << 4)
	// remaining length
	buffer.WriteByte(0x02)
	// packet id
	buffer.WriteByte(byte(p.PacketId >> 8))
	return buffer.WriteByte(byte(p.PacketId))
}
