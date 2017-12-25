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

type Packet interface {
	Type() CtrlType

	// Bytes dump a mqtt Packet object to mqtt bytes into provided buffer
	Bytes(*bytes.Buffer) error
}

type TopicMsg PublishPacket

type Topic struct {
	Name string
	Qos  QosLevel
}

const (
	maxMsgSize = 0xffffff7f
)

type CtrlType = byte

const (
	CtrlConn CtrlType = iota + 1
	CtrlConnAck
	CtrlPublish
	CtrlPubAck
	CtrlPubRecv
	CtrlPubRel
	CtrlPubComp
	CtrlSubscribe
	CtrlSubAck
	CtrlUnSub
	CtrlUnSubAck
	CtrlPingReq
	CtrlPingResp
	CtrlDisConn
)

type ProtocolLevel = byte

const (
	V31 ProtocolLevel = iota + 3
	V311
)

type QosLevel = byte

const (
	Qos0 QosLevel = iota
	Qos1
	Qos2
)

var (
	mqtt = []byte("MQTT")
)

type ConAckCode = byte

const (
	ConnAccepted ConAckCode = iota
	ConnBadProtocol
	ConnIdRejected
	ConnServerUnavailable
	ConnBadIdentity
	ConnAuthFail
	ConnTimeout   ConAckCode = 0xf0
	ConnBadPacket ConAckCode = 0xf1
	ConnDialErr   ConAckCode = 0xf2
)

type PubAckCode = byte

type SubAckCode = byte

const (
	SubOkMaxQos0 SubAckCode = iota
	SubOkMaxQos1
	SubOkMaxQos2
	SubFail = 0x80
)
