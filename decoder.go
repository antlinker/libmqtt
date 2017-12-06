package libmqtt

import (
	"bufio"
	"errors"
	"io"
)

var (
	ErrInvalidPacket = errors.New("read invalid packet ")
)

func decodeOnePacket(reader *bufio.Reader) (pkt Packet, err error) {
	var header byte
	header, err = reader.ReadByte()
	if err != nil {
		return
	}

	var bytesToRead int
	bytesToRead, err = decodeRemainLength(reader)
	if err != nil {
		return
	}

	body := make([]byte, bytesToRead)
	n, err := io.ReadFull(reader, body[:])
	if n < 2 {
		err = ErrInvalidPacket
	}
	if err != nil {
		return
	}

	switch header >> 4 {
	case CtrlConn:
		protocol, next := decodeString(body)
		if len(next) < 2 {
			err = ErrInvalidPacket
			return
		}
		hasUsername := next[1]&0x80>>7 == 1
		hasPassword := next[1]&0x40>>6 == 1
		tmpPkt := &ConnPacket{
			ProtoName:    protocol,
			ProtoLevel:   next[0],
			CleanSession: next[1]&0x02>>1 == 1,
			IsWill:       next[1]&0x04>>2 == 1,
			WillQos:      next[1] & 0x18 >> 3,
			WillRetain:   next[1]&0x20>>5 == 1,
			Keepalive:    uint16(next[2])<<8 + uint16(next[3]),
		}
		tmpPkt.ClientId, next = decodeString(next[4:])
		if tmpPkt.IsWill {
			tmpPkt.WillTopic, next = decodeString(next)
			tmpPkt.WillMessage, next = decodeString(next)
		}
		if hasUsername {
			tmpPkt.Username, next = decodeString(next)
		}
		if hasPassword {
			tmpPkt.Password, _ = decodeString(next)
		}
		pkt = tmpPkt
	case CtrlConnAck:
		if len(body) < 2 {
			err = ErrInvalidPacket
			return
		}
		pkt = &ConnAckPacket{
			Present: body[0]&0x01 == 1,
			Code:    body[1],
		}
	case CtrlPublish:
		topicName, next := decodeString(body)
		if len(next) < 2 {
			err = ErrInvalidPacket
			return
		}
		tmpPkg := &PublishPacket{
			IsDup:     header&0x08>>3 == 1,
			Qos:       header & 0x06 >> 1,
			IsRetain:  header&0x01 == 1,
			TopicName: topicName,
		}
		if tmpPkg.Qos > Qos0 {
			tmpPkg.PacketId = uint16(next[0])<<8 + uint16(next[1])
			next = next[2:]
		}
		tmpPkg.Payload = next[:]
		pkt = tmpPkg
	case CtrlPubAck:
		pkt = &PubAckPacket{
			PacketId: uint16(body[0])<<8 + uint16(body[1]),
		}
	case CtrlPubRecv:
		pkt = &PubRecvPacket{
			PacketId: uint16(body[0])<<8 + uint16(body[1]),
		}
	case CtrlPubRel:
		pkt = &PubRelPacket{
			PacketId: uint16(body[0])<<8 + uint16(body[1]),
		}
	case CtrlPubComp:
		pkt = &PubCompPacket{
			PacketId: uint16(body[0])<<8 + uint16(body[1]),
		}
	case CtrlSubscribe:
		pktTmp := &SubscribePacket{
			PacketId: uint16(body[0])<<8 + uint16(body[1]),
		}
		next := body[2:]
		topics := make([]*Topic, 0)
		for len(next) > 0 {
			var name string
			name, next = decodeString(next)
			topics = append(topics, &Topic{Name: name, Qos: next[0]})
			next = next[1:]
		}
		pktTmp.Topics = topics
		pkt = pktTmp
	case CtrlSubAck:
		pktTmp := &SubAckPacket{
			PacketId: uint16(body[0])<<8 + uint16(body[1]),
		}
		next := body[2:]
		codes := make([]SubAckCode, 0)
		for i := 0; i < len(next); i++ {
			codes = append(codes, next[i])
		}
		pkt = pktTmp
	case CtrlUnSub:
		pktTmp := &UnSubPacket{
			PacketId: uint16(body[0])<<8 + uint16(body[1]),
		}
		next := body[2:]
		topics := make([]string, 0)
		for len(next) > 0 {
			var name string
			name, next = decodeString(next)
			topics = append(topics, name)
		}
		pktTmp.TopicNames = topics
		pkt = pktTmp
	case CtrlUnSubAck:
		pkt = &UnSubAckPacket{
			PacketId: uint16(body[0])<<8 + uint16(body[1]),
		}
	case CtrlPingReq:
		pkt = &PingReqPacket{}
	case CtrlPingResp:
		pkt = &PingRespPacket{}
	case CtrlDisConn:
		pkt = &DisConnPacket{}
	}
	return
}

func decodeString(data []byte) (string, []byte) {
	if len(data) < 2 {
		return "", data
	}
	length := int(data[0])<<8 + int(data[1])
	return string(data[2 : 2+length]), data[2+length:]
}

func decodeRemainLength(reader io.ByteReader) (result int, err error) {
	m := 1
	var encodedByte byte
	encodedByte, err = reader.ReadByte()
	result = int(encodedByte & 127)
	for (encodedByte & 0x80) != 0 {
		encodedByte, err = reader.ReadByte()
		if err != nil {
			return
		}

		result += int(encodedByte&127) * m
		m <<= 8

		if m > 128*128*128 {
			return
		}
	}

	return
}
