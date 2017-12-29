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
#include <stdio.h>
#include <stdbool.h>
#include "libmqtth.h"

static inline void call_conn_handler(libmqtt_conn_handler h,
	char * server, libmqtt_connack_t code, char * err) {
  if (h != NULL) {
    h(server, code, err);
  }
}

static inline void call_pub_handler(libmqtt_pub_handler h,
	char * topic, char * msg) {
  if (h != NULL) {
    h(topic, msg);
  }
}

static inline void call_sub_handler(libmqtt_sub_handler h,
	char * topic, int qos, char * err) {
  if (h != NULL) {
    h(topic, qos, err);
  }
}

static inline void call_unsub_handler(libmqtt_unsub_handler h,
	char * topic, char * err) {
  if (h != NULL) {
    h(topic, err);
  }
}

static inline void call_net_handler(libmqtt_net_handler h,
	char * server, char * err) {
  if (h != NULL) {
    h(server, err);
  }
}

static inline void call_topic_handler(libmqtt_topic_handler h,
	char * topic, int qos , char * msg, int size) {
  if (h != NULL) {
    h(topic, qos, msg, size);
  }
}
*/
import "C"
import (
	"unsafe"

	mq "github.com/goiiot/libmqtt"
)

var (
	client  mq.Client
	options []mq.Option
)

// SetServer (char *server)
//export SetServer
func SetServer(server *C.char) {
	addOption(mq.WithServer(C.GoString(server)))
}

// SetCleanSession (bool flag)
//export SetCleanSession
func SetCleanSession(flag C.bool) {
	addOption(mq.WithCleanSession(bool(flag)))
}

// SetClientID (char * client_id)
//export SetClientID
func SetClientID(clientID *C.char) {
	addOption(mq.WithClientID(C.GoString(clientID)))
}

// SetDialTimeout (int timeout)
//export SetDialTimeout
func SetDialTimeout(timeout C.int) {
	addOption(mq.WithDialTimeout(uint16(timeout)))
}

// SetIdentity (char * username, char * password)
//export SetIdentity
func SetIdentity(username, password *C.char) {
	addOption(mq.WithIdentity(C.GoString(username), C.GoString(password)))
}

// SetKeepalive (int keepalive, float factor)
//export SetKeepalive
func SetKeepalive(keepalive C.int, factor C.float) {
	addOption(mq.WithKeepalive(uint16(keepalive), float64(factor)))
}

// SetLog (libmqtt_log_level l)
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

// SetRecvBuf (int size)
//export SetRecvBuf
func SetRecvBuf(size C.int) {
	addOption(mq.WithRecvBuf(int(size)))
}

// SetSendBuf (int size)
//export SetSendBuf
func SetSendBuf(size C.int) {
	addOption(mq.WithSendBuf(int(size)))
}

// SetTLS (char * certFile, char * keyFile, char * caCert, char * serverNameOverride, bool skipVerify)
// use ssl to connect
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

// SetWill (char *topic, int qos, bool retain, char *payload, int payloadSize)
// mark this connection with will message
//export SetWill
func SetWill(topic *C.char, qos C.int, retain C.bool, payload *C.char, payloadSize C.int) {
	addOption(mq.WithWill(
		C.GoString(topic),
		mq.QosLevel(qos),
		bool(retain),
		C.GoBytes(unsafe.Pointer(payload), payloadSize),
	))
}

// SetUp ()
// setup the client with previously defined options
//export SetUp
func SetUp() *C.char {
	var err error
	client, err = mq.NewClient(options...)
	if err != nil {
		return C.CString(err.Error())
	}

	client.HandlePub(pubHandler)
	client.HandleSub(subHandler)
	client.HandleUnSub(unSubHandler)
	client.HandleNet(netHandler)

	return nil
}

// Handle (char *topic, libmqtt_topic_handler h)
//export Handle
func Handle(topic *C.char, h C.libmqtt_topic_handler) {
	if client != nil {
		client.Handle(C.GoString(topic), wrapTopicHandler(h))
	}
}

// Conn ()
//export Conn
func Conn() {
	if client != nil {
		client.Connect(connHandler)
	}
}

// Subscribe (char *topic, int qos)
//export Subscribe
func Subscribe(topic *C.char, qos C.int) {
	if client != nil {
		client.Subscribe(&mq.Topic{
			Name: C.GoString(topic),
			Qos:  mq.QosLevel(qos),
		})
	}
}

// Publish (char *topic, int qos, char *payload, int payloadSize)
//export Publish
func Publish(topic *C.char, qos C.int, payload *C.char, payloadSize C.int) {
	if client != nil {
		client.Publish(&mq.PublishPacket{
			TopicName: C.GoString(topic),
			Qos:       mq.QosLevel(qos),
			Payload:   C.GoBytes(unsafe.Pointer(payload), payloadSize),
		})
	}
}

// UnSubscribe (char *topic)
//export UnSubscribe
func UnSubscribe(topic *C.char) {
	if client != nil {
		client.UnSubscribe(C.GoString(topic))
	}
}

// Wait ()
//export Wait
func Wait() {
	if client != nil {
		client.Wait()
	}
}

// Destroy (bool force)
//export Destroy
func Destroy(force C.bool) {
	if client != nil {
		client.Destroy(bool(force))
	}
}

func main() {}

func addOption(o ...mq.Option) {
	if options == nil {
		options = make([]mq.Option, 0)
	}
	options = append(options, o...)
}

func connHandler(server string, code mq.ConnAckCode, err error) {
	if mConnHandler != nil {
		var c C.libmqtt_connack_t
		switch code {
		case mq.ConnAccepted:
			c = C.libmqtt_connack_accepted
		case mq.ConnBadProtocol:
			c = C.libmqtt_connack_bad_proto
		case mq.ConnIDRejected:
			c = C.libmqtt_connack_id_rejected
		case mq.ConnServerUnavailable:
			c = C.libmqtt_connack_srv_unavail
		case mq.ConnBadIdentity:
			c = C.libmqtt_connack_bad_identity
		case mq.ConnAuthFail:
			c = C.libmqtt_connack_auth_fail
		}

		var er *C.char
		if err != nil {
			er = C.CString(err.Error())
		}
		C.call_conn_handler(mConnHandler, C.CString(server), c, er)
	}
}

func pubHandler(topic string, err error) {
	if mPubHandler != nil {
		var er *C.char
		if err != nil {
			er = C.CString(err.Error())
		}
		C.call_pub_handler(mPubHandler, C.CString(topic), er)
	}
}

func subHandler(topics []*mq.Topic, err error) {
	if mSubHandler != nil {
		for _, t := range topics {
			var er *C.char
			if err != nil {
				er = C.CString(err.Error())
			}
			C.call_sub_handler(mSubHandler, C.CString(t.Name), C.int(t.Qos), er)
		}
	}
}

func unSubHandler(topics []string, err error) {
	if mSubHandler != nil {
		for _, t := range topics {
			var er *C.char
			if err != nil {
				er = C.CString(err.Error())
			}
			C.call_unsub_handler(mUnSubHandler, C.CString(t), er)
		}
	}
}

func netHandler(server string, err error) {
	if mNetHandler != nil {
		var er *C.char
		if err != nil {
			er = C.CString(err.Error())
		}
		C.call_net_handler(mNetHandler, C.CString(server), er)
	}
}

func wrapTopicHandler(h C.libmqtt_topic_handler) mq.TopicHandler {
	return func(topic string, qos mq.QosLevel, msg []byte) {
		C.call_topic_handler(h, C.CString(topic), C.int(qos), (*C.char)(C.CBytes(msg)), C.int(len(msg)))
	}
}
