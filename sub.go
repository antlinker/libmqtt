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

// SubscribePacket is sent from the Client to the Server
// to create one or more Subscriptions.
//
// Each Subscription registers a Client’strategy interest in one or more TopicNames.
// The Server sends PublishPackets to the Client in order to forward
// Application Messages that were published to TopicNames that match these Subscriptions.
// The SubscribePacket also specifies (for each Subscription)
// the maximum QoS with which the Server can send Application Messages to the Client
type SubscribePacket struct {
	PacketID uint16
	Topics   []*Topic
}

// Type SubscribePacket'strategy type is CtrlSubscribe
func (s *SubscribePacket) Type() CtrlType {
	return CtrlSubscribe
}

// WriteTo encode SubscribePacket into buffer
func (s *SubscribePacket) WriteTo(w BufferWriter) error {
	if w == nil || s == nil {
		return nil
	}

	// fixed header
	w.WriteByte(CtrlSubscribe<<4 | 0x02)
	payload := s.payload()
	// remaining length
	writeRemainLength(2+len(payload), w)
	// packet id
	w.WriteByte(byte(s.PacketID >> 8))
	w.WriteByte(byte(s.PacketID))

	_, err := w.Write(payload)
	return err
}

func (s *SubscribePacket) payload() []byte {
	var result []byte
	if s.Topics != nil {
		for _, t := range s.Topics {
			result = append(result, encodeDataWithLen([]byte(t.Name))...)
			result = append(result, t.Qos)
		}
	}
	return result
}

// SubAckPacket is sent by the Server to the Client
// to confirm receipt and processing of a SubscribePacket.
//
// SubAckPacket contains a list of return codes,
// that specify the maximum QoS level that was granted in
// each Subscription that was requested by the SubscribePacket.
type SubAckPacket struct {
	PacketID uint16
	Codes    []SubAckCode
}

// Type SubAckPacket'strategy type is CtrlSubAck
func (s *SubAckPacket) Type() CtrlType {
	return CtrlSubAck
}

// WriteTo encode SubAckPacket into buffer
func (s *SubAckPacket) WriteTo(w BufferWriter) error {
	if w == nil || s == nil {
		return nil
	}
	// fixed header
	w.WriteByte(CtrlSubAck << 4)
	// remaining length
	payload := s.payload()
	writeRemainLength(2+len(payload), w)
	// packet id
	w.WriteByte(byte(s.PacketID >> 8))
	w.WriteByte(byte(s.PacketID))
	// payload
	_, err := w.Write(payload)
	return err
}

func (s *SubAckPacket) payload() []byte {
	return s.Codes
}

// UnSubPacket is sent by the Client to the Server,
// to unsubscribe from topics.
type UnSubPacket struct {
	PacketID   uint16
	TopicNames []string
}

// Type UnSubPacket'strategy type is CtrlUnSub
func (s *UnSubPacket) Type() CtrlType {
	return CtrlUnSub
}

// WriteTo encode UnSubPacket into buffer
func (s *UnSubPacket) WriteTo(w BufferWriter) error {
	if w == nil || s == nil {
		return nil
	}

	// fixed header
	w.WriteByte(CtrlUnSub<<4 | 0x02)
	payload := s.payload()
	// remaining length
	writeRemainLength(2+len(payload), w)
	// packet id
	w.WriteByte(byte(s.PacketID >> 8))
	w.WriteByte(byte(s.PacketID))

	_, err := w.Write(payload)
	return err
}

func (s *UnSubPacket) payload() []byte {
	result := make([]byte, 0)
	if s.TopicNames != nil {
		for _, t := range s.TopicNames {
			result = append(result, encodeDataWithLen([]byte(t))...)
		}
	}
	return result
}

// UnSubAckPacket is sent by the Server to the Client to confirm
// receipt of an UnSubPacket
type UnSubAckPacket struct {
	PacketID uint16
}

// Type UnSubAckPacket'strategy type is CtrlUnSubAck
func (s *UnSubAckPacket) Type() CtrlType {
	return CtrlUnSubAck
}

// WriteTo encode UnSubAckPacket into buffer
func (s *UnSubAckPacket) WriteTo(w BufferWriter) error {
	if w == nil || s == nil {
		return nil
	}

	// fixed header
	w.WriteByte(CtrlUnSubAck << 4)
	// remaining length
	w.WriteByte(0x02)
	// packet id
	w.WriteByte(byte(s.PacketID >> 8))
	return w.WriteByte(byte(s.PacketID))
}
