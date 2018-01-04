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

// ConPacket is the first packet sent by Client to Server
type ConPacket struct {
	protoName    string
	protoLevel   ProtocolLevel
	Username     string
	Password     string
	ClientID     string
	CleanSession bool
	IsWill       bool
	WillQos      QosLevel
	WillRetain   bool
	Keepalive    uint16
	WillTopic    string
	WillMessage  []byte
}

// Type ConPacket'strategy type is CtrlConn
func (c *ConPacket) Type() CtrlType {
	return CtrlConn
}

// Bytes encode ConPacket to bytes
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
	var connectFlag byte
	if c.ClientID == "" {
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
	encodeDataWithLen([]byte(c.ClientID), result)

	// will topic and message
	if c.IsWill {
		encodeDataWithLen([]byte(c.WillTopic), result)
		encodeDataWithLen(c.WillMessage, result)
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
	Code    ConnAckCode
}

// Type ConAckPacket'strategy type is CtrlConnAck
func (c *ConAckPacket) Type() CtrlType {
	return CtrlConnAck
}

// Bytes encode ConAckPacket to bytes
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

var (
	// DisConPacket is the final instance of disConPacket
	DisConPacket = &disConPacket{}
)

// disConPacket is the final Control Packet sent from the Client to the Server.
// It indicates that the Client is disconnecting cleanly.
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
