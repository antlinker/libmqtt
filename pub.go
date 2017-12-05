package libmqtt

type PublishPacket struct {
	Dup      bool
	Qos      QosLevel
	Retain   bool
	Topic    string
	PacketId uint16
	Payload  []byte
}

func (p *PublishPacket) Type() CtrlType {
	return CtrlPublish
}

func (p *PublishPacket) flag() byte {
	return 0x00
}

// PubAckPacket is the response to a PublishPacket with QoS level 1.
type PubAckPacket struct {
	PacketId uint16
}

func (p *PubAckPacket) Type() CtrlType {
	return CtrlPubAck
}

func (p *PubAckPacket) flag() byte {
	return 0x00
}

// PubRecPacket is the response to a PublishPacket with QoS 2.
// It is the second packet of the QoS 2 protocol exchange.
type PubRecPacket struct {
	PacketId uint16
}

func (p *PubRecPacket) Type() CtrlType {
	return CtrlPubRecv
}

func (p *PubRecPacket) flag() byte {
	return 0x00
}

// PubRelPacket is the response to a PubRecPacket.
// It is the third packet of the QoS 2 protocol exchange.
type PubRelPacket struct {
	PacketId uint16
}

func (p *PubRelPacket) Type() CtrlType {
	return CtrlPubRel
}

func (p *PubRelPacket) flag() byte {
	return 0x00
}

// PubCompPacket is the response to a PubRelPacket.
// It is the fourth and final packet of the QoS 892 2 protocol exchange. 893
type PubCompPacket struct {
	PacketId uint16
}

func (p *PubCompPacket) Type() CtrlType {
	return CtrlPubComp
}

func (p *PubCompPacket) flag() byte {
	return 0x00
}
