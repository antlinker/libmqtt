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

import (
	"bytes"
)

// Packet is MQTT control packet
type Packet interface {
	// Type return the packet type
	Type() CtrlType

	// Bytes dump a mqtt Packet object to mqtt bytes into provided buffer
	Bytes(*bytes.Buffer) error
}

// Topic for both topic name and topic qos
type Topic struct {
	Name string
	Qos  QosLevel
}

// String
func (t *Topic) String() string {
	return t.Name
}

const (
	maxMsgSize = 0xffffff7f
)

// CtrlType is MQTT Control packet type
type CtrlType = byte

const (
	// CtrlConn Connect
	CtrlConn CtrlType = iota + 1
	// CtrlConnAck Connect Ack
	CtrlConnAck
	// CtrlPublish Publish
	CtrlPublish
	// CtrlPubAck Publish Ack
	CtrlPubAck
	// CtrlPubRecv Publish Received
	CtrlPubRecv
	// CtrlPubRel Publish Release
	CtrlPubRel
	// CtrlPubComp Publish Complete
	CtrlPubComp
	// CtrlSubscribe Subscribe
	CtrlSubscribe
	// CtrlSubAck Subscribe Ack
	CtrlSubAck
	// CtrlUnSub UnSubscribe
	CtrlUnSub
	// CtrlUnSubAck UnSubscribe Ack
	CtrlUnSubAck
	// CtrlPingReq Ping Request
	CtrlPingReq
	// CtrlPingResp Ping Response
	CtrlPingResp
	// CtrlDisConn Disconnect
	CtrlDisConn
)

// ProtocolLevel MQTT Protocol
type ProtocolLevel = byte

// V311 means MQTT 3.1.1
const V311 ProtocolLevel = 4

// QosLevel is either 0, 1, 2
type QosLevel = byte

const (
	// Qos0 0
	Qos0 QosLevel = iota
	// Qos1 1
	Qos1
	// Qos2 2
	Qos2
)

var (
	mqtt = []byte("MQTT")
)

// ConnAckCode is connection response code from server
type ConnAckCode = byte

const (
	// ConnAccepted client accepted by server
	ConnAccepted ConnAckCode = iota
	// ConnBadProtocol Protocol not supported
	ConnBadProtocol
	// ConnIDRejected Connection Id not valid
	ConnIDRejected
	// ConnServerUnavailable Server error
	ConnServerUnavailable
	// ConnBadIdentity Identity failed
	ConnBadIdentity
	// ConnAuthFail Auth failed
	ConnAuthFail
)

// SubAckCode is returned by server in SubAckPacket
type SubAckCode = byte

const (
	// SubOkMaxQos0 QoS 0 is used by server
	SubOkMaxQos0 SubAckCode = iota
	// SubOkMaxQos1 QoS 1 is used by server
	SubOkMaxQos1
	// SubOkMaxQos2 QoS 2 is used by server
	SubOkMaxQos2
	// SubFail means that subscription is not successful
	SubFail SubAckCode = 0x80
)
