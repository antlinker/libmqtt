/*
 * Copyright GoIIoT (https://github.com/goiiot)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
	PacketID  uint16
}

// Type PublishPacket's type is CtrlPublish
func (p *PublishPacket) Type() CtrlType {
	return CtrlPublish
}

// Bytes encode PublishPacket into buffer
func (p *PublishPacket) Bytes(buffer *bytes.Buffer) (err error) {
	if buffer == nil || p == nil {
		return
	}
	// fixed header
	buffer.WriteByte(CtrlPublish<<4 | boolToByte(p.IsDup)<<3 |
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
		result.WriteByte(byte(p.PacketID >> 8))
		result.WriteByte(byte(p.PacketID))
	}
	result.Write(p.Payload)
	return
}

// PubAckPacket is the response to a PublishPacket with QoS level 1.
type PubAckPacket struct {
	PacketID uint16
}

// Type PubAckPacket's type is CtrlPubAck
func (p *PubAckPacket) Type() CtrlType {
	return CtrlPubAck
}

// Bytes encode PubAckPacket into buffer
func (p *PubAckPacket) Bytes(buffer *bytes.Buffer) (err error) {
	if buffer == nil || p == nil {
		return
	}

	// fixed header
	buffer.WriteByte(CtrlPubAck << 4)
	// remaining length
	buffer.WriteByte(0x02)
	// packet id
	buffer.WriteByte(byte(p.PacketID >> 8))
	return buffer.WriteByte(byte(p.PacketID))
}

// PubRecvPacket is the response to a PublishPacket with QoS 2.
// It is the second packet of the QoS 2 protocol exchange.
type PubRecvPacket struct {
	PacketID uint16
}

// Type PubRecvPacket's type is CtrlPubRecv
func (p *PubRecvPacket) Type() CtrlType {
	return CtrlPubRecv
}

// Bytes encode PubRecvPacket into buffer
func (p *PubRecvPacket) Bytes(buffer *bytes.Buffer) (err error) {
	if buffer == nil || p == nil {
		return
	}

	// fixed header
	buffer.WriteByte(CtrlPubRecv << 4)
	// remaining length
	buffer.WriteByte(0x02)
	// packet id
	buffer.WriteByte(byte(p.PacketID >> 8))
	return buffer.WriteByte(byte(p.PacketID))
}

// PubRelPacket is the response to a PubRecvPacket.
// It is the third packet of the QoS 2 protocol exchange.
type PubRelPacket struct {
	PacketID uint16
}

// Type PubRelPacket's type is CtrlPubRel
func (p *PubRelPacket) Type() CtrlType {
	return CtrlPubRel
}

// Bytes encode PubRelPacket into buffer
func (p *PubRelPacket) Bytes(buffer *bytes.Buffer) (err error) {
	if buffer == nil || p == nil {
		return
	}

	buffer.WriteByte(CtrlPubRel<<4 | 0x02)
	// remaining length
	buffer.WriteByte(0x02)
	// packet id
	buffer.WriteByte(byte(p.PacketID >> 8))
	return buffer.WriteByte(byte(p.PacketID))
}

// PubCompPacket is the response to a PubRelPacket.
// It is the fourth and final packet of the QoS 892 2 protocol exchange. 893
type PubCompPacket struct {
	PacketID uint16
}

// Type PubCompPacket's type is CtrlPubComp
func (p *PubCompPacket) Type() CtrlType {
	return CtrlPubComp
}

// Bytes encode PubCompPacket into buffer
func (p *PubCompPacket) Bytes(buffer *bytes.Buffer) (err error) {
	if buffer == nil || p == nil {
		return
	}
	// fixed header
	buffer.WriteByte(CtrlPubComp << 4)
	// remaining length
	buffer.WriteByte(0x02)
	// packet id
	buffer.WriteByte(byte(p.PacketID >> 8))
	return buffer.WriteByte(byte(p.PacketID))
}
