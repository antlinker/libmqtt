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

// SubscribePacket is sent from the Client to the Server
// to create one or more Subscriptions.
//
// Each Subscription registers a Client’s interest in one or more TopicNames.
// The Server sends PublishPackets to the Client in order to forward
// Application Messages that were published to TopicNames that match these Subscriptions.
// The SubscribePacket also specifies (for each Subscription)
// the maximum QoS with which the Server can send Application Messages to the Client
type SubscribePacket struct {
	PacketId uint16
	Topics   []*Topic
}

func (s *SubscribePacket) Type() CtrlType {
	return CtrlSubscribe
}

func (s *SubscribePacket) Bytes(buffer *bytes.Buffer) (err error) {
	if buffer == nil || s == nil {
		return
	}

	// fixed header
	buffer.WriteByte((CtrlSubscribe << 4) | 0x02)
	payload := s.payload()
	// remaining length
	encodeRemainLength(2+payload.Len(), buffer)
	// packet id
	buffer.WriteByte(byte(s.PacketId >> 8))
	buffer.WriteByte(byte(s.PacketId))

	_, err = payload.WriteTo(buffer)

	return
}

func (s *SubscribePacket) payload() (result *bytes.Buffer) {
	result = &bytes.Buffer{}
	if s.Topics != nil {
		for _, t := range s.Topics {
			lenTopicName := len(t.Name)
			result.WriteByte(byte(lenTopicName >> 8))
			result.WriteByte(byte(lenTopicName))
			result.Write([]byte(t.Name))
			result.WriteByte(t.Qos)
		}
	}
	return
}

// SubAckPacket is sent by the Server to the Client
// to confirm receipt and processing of a SubscribePacket.
//
// SubAckPacket contains a list of return codes,
// that specify the maximum QoS level that was granted in
// each Subscription that was requested by the SubscribePacket.
type SubAckPacket struct {
	PacketId uint16
	Codes    []SubAckCode
}

func (s *SubAckPacket) Type() CtrlType {
	return CtrlSubAck
}

func (s *SubAckPacket) Bytes(buffer *bytes.Buffer) (err error) {
	if buffer == nil || s == nil {
		return
	}
	// fixed header
	buffer.WriteByte(CtrlSubAck << 4)
	// remaining length
	payload := s.payload()
	encodeRemainLength(2+payload.Len(), buffer)
	// packet id
	buffer.WriteByte(byte(s.PacketId >> 8))
	buffer.WriteByte(byte(s.PacketId))
	// payload
	_, err = payload.WriteTo(buffer)

	return
}

func (s *SubAckPacket) payload() (result *bytes.Buffer) {
	result = &bytes.Buffer{}
	if s.Codes != nil {
		for _, c := range s.Codes {
			result.WriteByte(c)
		}
	}
	return
}

// UnSubPacket is sent by the Client to the Server,
// to unsubscribe from topics.
type UnSubPacket struct {
	PacketId   uint16
	TopicNames []string
}

func (s *UnSubPacket) Type() CtrlType {
	return CtrlUnSub
}

func (s *UnSubPacket) Bytes(buffer *bytes.Buffer) (err error) {
	if buffer == nil || s == nil {
		return
	}

	// fixed header
	buffer.WriteByte(CtrlUnSub << 4)
	payload := s.payload()
	// remaining length
	encodeRemainLength(2+payload.Len(), buffer)
	// packet id
	buffer.WriteByte(byte(s.PacketId >> 8))
	buffer.WriteByte(byte(s.PacketId))

	_, err = payload.WriteTo(buffer)

	return
}

func (s *UnSubPacket) payload() (result *bytes.Buffer) {
	result = &bytes.Buffer{}
	if s.TopicNames != nil {
		for _, t := range s.TopicNames {
			encodeDataWithLen([]byte(t), result)
		}
	}

	return
}

// UnSubAckPacket is sent by the Server to the Client to confirm
// receipt of an UnSubPacket
type UnSubAckPacket struct {
	PacketId uint16
}

func (s *UnSubAckPacket) Type() CtrlType {
	return CtrlUnSubAck
}

func (s *UnSubAckPacket) Bytes(buffer *bytes.Buffer) (err error) {
	if buffer == nil || s == nil {
		return
	}

	// fixed header
	buffer.WriteByte(CtrlUnSubAck << 4)
	// remaining length
	buffer.WriteByte(0x02)
	// packet id
	buffer.WriteByte(byte(s.PacketId >> 8))
	err = buffer.WriteByte(byte(s.PacketId))

	return
}
