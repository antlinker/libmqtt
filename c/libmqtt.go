// +build cgo lib

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

package main

// #cgo CFLAGS: -I include
/*
#include <string.h>
#include <stdbool.h>
#include "libmqtth.h"
*/
import "C"
import (
	mq "github.com/goiiot/libmqtt"
)

var (
	client  mq.Client
	options []mq.Option
)

//export SetServer
func SetServer(server *C.char) {
	addOption(mq.WithServer(C.GoString(server)))
}

//export SetCleanSession
func SetCleanSession(flag C.bool) {
	addOption(mq.WithCleanSession(bool(flag)))
}

//export SetClientID
func SetClientID(clientId *C.char) {
	addOption(mq.WithClientID(C.GoString(clientId)))
}

//export SetDialTimeout
func SetDialTimeout(timeout C.int) {
	addOption(mq.WithDialTimeout(uint16(timeout)))
}

//export SetIdentity
func SetIdentity(username, password *C.char) {
	addOption(mq.WithIdentity(C.GoString(username), C.GoString(password)))
}

//export SetKeepalive
func SetKeepalive(keepalive C.int, factor C.float) {
	addOption(mq.WithKeepalive(uint16(keepalive), float64(factor)))
}

//export SetLog
func SetLog(l C.libmqtt_log_level) {
	level := mq.Silent
	switch l {
	case C.libmqtt_log_verbose:
		level = mq.Verbose
	case C.libmqtt_log_debug:
		level = mq.Debug
	case C.libmqtt_log_info:
		level = mq.Info
	case C.libmqtt_log_warning:
		level = mq.Warning
	case C.libmqtt_log_error:
		level = mq.Error
	}
	addOption(mq.WithLog(level))
}

//export SetRecvBuf
func SetRecvBuf(size C.int) {
	addOption(mq.WithRecvBuf(int(size)))
}

//export SetSendBuf
func SetSendBuf(size C.int) {
	addOption(mq.WithSendBuf(int(size)))
}

//export SetTLS
func SetTLS(certFile, keyFile, caCert, serverNameOverride *C.char, skipVerify C.bool) {
	addOption(mq.WithTLS(
		C.GoString(certFile),
		C.GoString(keyFile),
		C.GoString(caCert),
		C.GoString(serverNameOverride),
		bool(skipVerify)),
	)
}

//export SetWill
func SetWill(topic *C.char, qos C.int, retain C.bool, payload *C.char, payloadSize C.size_t) {
	addOption(mq.WithWill(
		C.GoString(topic),
		mq.QosLevel(qos),
		bool(retain),
		goByteSlice(payload, payloadSize),
	))
}

//export SetUp
func SetUp() {
	client = mq.NewClient(options...)

	client.HandlePub(callPubHandler)
	client.HandleSub(callSubHandler)
	client.HandleUnSub(callUnSubHandler)
	client.HandleNet(callNetHandler)
}

//export Conn
func Conn() {
	if client != nil {
		client.Connect(callConnHandler)
	}
}

// Sub(char *topic, int qos)
//export Sub
func Sub(topic *C.char, qos C.int) {
	if client != nil {
		client.Subscribe(&mq.Topic{
			Name: C.GoString(topic),
			Qos:  mq.QosLevel(qos),
		})
	}
}

// Pub(char *topic, int qos, char *payload, size_t payloadSize)
//export Pub
func Pub(topic *C.char, qos C.int, payload *C.char, payloadSize C.size_t) {
	if client != nil {
		client.Publish(&mq.PublishPacket{
			TopicName: C.GoString(topic),
			Qos:       mq.QosLevel(qos),
			Payload:   goByteSlice(payload, payloadSize),
		})
	}
}

// UnSub(char *topic)
//export UnSub
func UnSub(topic *C.char) {
	if client != nil {
		client.UnSubscribe(C.GoString(topic))
	}
}

// Destroy(bool force)
//export Destroy
func Destroy(force C.bool) {
	if client != nil {
		client.Destroy(bool(force))
	}
}

func addOption(o ...mq.Option) {
	if options == nil {
		options = make([]mq.Option, 0)
	}
	options = append(options, o...)
}

func main() {}
