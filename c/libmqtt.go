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

typedef enum {
  libmqtt_connack_accepted = 0,
  libmqtt_connack_bad_proto = 1,
  libmqtt_connack_id_rejected = 2,
  libmqtt_connack_srv_unavail = 3,
  libmqtt_connack_bad_identity = 4,
  libmqtt_connack_auth_fail = 5,
} libmqtt_connack_t;

typedef enum {
  libmqtt_suback_max_qos0 = 0,
  libmqtt_suback_max_qos1 = 1,
  libmqtt_suback_max_qos2 = 2,
  libmqtt_suback_fail = 0x80,
} libmqtt_suback_t;

typedef enum {
  libmqtt_log_verbose = 0,
  libmqtt_log_debug = 1,
  libmqtt_log_info = 2,
  libmqtt_log_warning = 3,
  libmqtt_log_error = 4,
  libmqtt_log_silent = 5,
} libmqtt_log_level;

typedef void (*libmqtt_conn_handler)
  (char *server, libmqtt_connack_t code, char *err);

typedef void (*libmqtt_pub_handler)
  (char *topic, char *err);

typedef void (*libmqtt_sub_handler)
  (char *topic, int qos, char *err);

typedef void (*libmqtt_unsub_handler)
  (char *topic, char *err);

typedef void (*libmqtt_net_handler)
  (char *server, char *err);

typedef void (*libmqtt_persist_handler)
  (char *err);

typedef void (*libmqtt_topic_handler)
  (char *topic, int qos, char *msg, int size);

static inline void call_conn_handler
  (libmqtt_conn_handler h, char * server, libmqtt_connack_t code, char * err) {
  if (h != NULL) {
    h(server, code, err);
  }
}

static inline void call_pub_handler
  (libmqtt_pub_handler h, char * topic, char * msg) {
  if (h != NULL) {
    h(topic, msg);
  }
}

static inline void call_sub_handler
  (libmqtt_sub_handler h, char * topic, int qos, char * err) {
  if (h != NULL) {
    h(topic, qos, err);
  }
}

static inline void call_unsub_handler
  (libmqtt_unsub_handler h, char * topic, char * err) {
  if (h != NULL) {
    h(topic, err);
  }
}

static inline void call_net_handler
  (libmqtt_net_handler h, char * server, char * err) {
  if (h != NULL) {
    h(server, err);
  }
}

static inline void call_persist_handler
  (libmqtt_persist_handler h, char *err) {
  if (h != NULL) {
    h(err);
  }
}

static inline void call_topic_handler
  (libmqtt_topic_handler h, char * topic, int qos , char * msg, int size) {
  if (h != NULL) {
    h(topic, qos, msg, size);
  }
}
*/
import "C"
import (
	"math"
	"unsafe"

	mq "github.com/goiiot/libmqtt"
)

var (
	clients       = make(map[int]mq.Client)
	clientOptions = make(map[int][]mq.Option)
)

// Libmqtt_new_client ()
// create a new client, return the id of client
// if no client id available return -1
//export Libmqtt_new_client
func Libmqtt_new_client() C.int {
	for id := 1; id < int(math.MaxInt64); id++ {
		if _, ok := clientOptions[id]; !ok {
			clientOptions[id] = make([]mq.Option, 0)
			return C.int(id)
		}
	}
	return -1
}

// Libmqtt_client_set_server (int id, char *server)
//export Libmqtt_client_set_server
func Libmqtt_client_set_server(id C.int, server *C.char) {
	cid := int(id)
	if v, ok := clientOptions[cid]; ok {
		v = append(v, mq.WithServer(C.GoString(server)))
		clientOptions[cid] = v
	}
}

// Libmqtt_client_set_clean_session (bool flag)
//export Libmqtt_client_set_clean_session
func Libmqtt_client_set_clean_session(id C.int, flag C.bool) {
	cid := int(id)
	if v, ok := clientOptions[cid]; ok {
		v = append(v, mq.WithCleanSession(bool(flag)))
		clientOptions[cid] = v
	}
}

// Libmqtt_client_set_client_id (int id, char * client_id)
//export Libmqtt_client_set_client_id
func Libmqtt_client_set_client_id(id C.int, clientID *C.char) {
	cid := int(id)
	if v, ok := clientOptions[cid]; ok {
		v = append(v, mq.WithClientID(C.GoString(clientID)))
		clientOptions[cid] = v
	}
}

// Libmqtt_client_set_dial_timeout (int id, int timeout)
//export Libmqtt_client_set_dial_timeout
func Libmqtt_client_set_dial_timeout(id C.int, timeout C.int) {
	cid := int(id)
	if v, ok := clientOptions[cid]; ok {
		v = append(v, mq.WithDialTimeout(uint16(timeout)))
		clientOptions[cid] = v
	}
}

// Libmqtt_client_set_identity (int id, char * username, char * password)
//export Libmqtt_client_set_identity
func Libmqtt_client_set_identity(id C.int, username, password *C.char) {
	cid := int(id)
	if v, ok := clientOptions[cid]; ok {
		v = append(v, mq.WithIdentity(C.GoString(username), C.GoString(password)))
		clientOptions[cid] = v
	}
}

// Libmqtt_client_set_keepalive (int id, int keepalive, float factor)
//export Libmqtt_client_set_keepalive
func Libmqtt_client_set_keepalive(id C.int, keepalive C.int, factor C.float) {
	cid := int(id)
	if v, ok := clientOptions[cid]; ok {
		v = append(v, mq.WithKeepalive(uint16(keepalive), float64(factor)))
		clientOptions[cid] = v
	}
}

// Libmqtt_client_set_log (int id, libmqtt_log_level l)
//export Libmqtt_client_set_log
func Libmqtt_client_set_log(id C.int, l C.libmqtt_log_level) {
	cid := int(id)
	if v, ok := clientOptions[cid]; ok {
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
		case C.libmqtt_log_silent:
			level = mq.Silent
		}
		v = append(v, mq.WithLog(level))
		clientOptions[cid] = v
	}
}

// Libmqtt_client_set_recv_buf (int id, int size)
//export Libmqtt_client_set_recv_buf
func Libmqtt_client_set_recv_buf(id C.int, size C.int) {
	cid := int(id)
	if v, ok := clientOptions[cid]; ok {
		v = append(v, mq.WithRecvBuf(int(size)))
		clientOptions[cid] = v
	}
}

// Libmqtt_client_set_send_buf (int id, int size)
//export Libmqtt_client_set_send_buf
func Libmqtt_client_set_send_buf(id C.int, size C.int) {
	cid := int(id)
	if v, ok := clientOptions[cid]; ok {
		v = append(v, mq.WithSendBuf(int(size)))
		clientOptions[cid] = v
	}
}

// Libmqtt_client_set_tls (int id, char * certFile, char * keyFile, char * caCert, char * serverNameOverride, bool skipVerify)
// use ssl to connect
//export Libmqtt_client_set_tls
func Libmqtt_client_set_tls(id C.int, certFile, keyFile, caCert, serverNameOverride *C.char, skipVerify C.bool) {
	cid := int(id)
	if v, ok := clientOptions[cid]; ok {
		v = append(v, mq.WithTLS(
			C.GoString(certFile),
			C.GoString(keyFile),
			C.GoString(caCert),
			C.GoString(serverNameOverride),
			bool(skipVerify)),
		)
		clientOptions[cid] = v
	}
}

// Libmqtt_client_set_will (int id, char *topic, int qos, bool retain, char *payload, int payloadSize)
// mark this connection with will message
//export Libmqtt_client_set_will
func Libmqtt_client_set_will(id C.int, topic *C.char, qos C.int, retain C.bool, payload *C.char, payloadSize C.int) {
	cid := int(id)
	if v, ok := clientOptions[cid]; ok {
		v = append(v, mq.WithWill(
			C.GoString(topic),
			mq.QosLevel(qos),
			bool(retain),
			C.GoBytes(unsafe.Pointer(payload), payloadSize)))
		clientOptions[cid] = v
	}
}

// Libmqtt_setup (int id)
// setup the client with previously defined options
//export Libmqtt_setup
func Libmqtt_setup(id C.int) *C.char {
	cid := int(id)
	if v, ok := clientOptions[cid]; ok {
		c, err := mq.NewClient(v...)
		if err != nil {
			return C.CString(err.Error())
		}

		clients[cid] = c
	} else {
		return C.CString("no such client")
	}

	return nil
}

// Libmqtt_handle (int id, char *topic, libmqtt_topic_handler h)
//export Libmqtt_handle
func Libmqtt_handle(id C.int, topic *C.char, h C.libmqtt_topic_handler) {
	if c, ok := clients[int(id)]; ok {
		c.Handle(C.GoString(topic), wrapTopicHandler(h))
	}
}

// Libmqtt_connect (int id)
//export Libmqtt_connect
func Libmqtt_connect(id C.int, h C.libmqtt_conn_handler) {
	cid := int(id)
	if c, ok := clients[cid]; ok {
		c.Connect(func(server string, code mq.ConnAckCode, err error) {
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
			C.call_conn_handler(h, C.CString(server), c, er)
		})
	}
}

// Libmqtt_subscribe (int id, char *topic, int qos)
//export Libmqtt_subscribe
func Libmqtt_subscribe(id C.int, topic *C.char, qos C.int) {
	cid := int(id)
	if c, ok := clients[cid]; ok {
		c.Subscribe(&mq.Topic{
			Name: C.GoString(topic),
			Qos:  mq.QosLevel(qos),
		})
	}
}

// Libmqtt_publish (int id, char *topic, int qos, char *payload, int payloadSize)
//export Libmqtt_publish
func Libmqtt_publish(id C.int, topic *C.char, qos C.int, payload *C.char, payloadSize C.int) {
	cid := int(id)
	if c, ok := clients[cid]; ok {
		c.Publish(&mq.PublishPacket{
			TopicName: C.GoString(topic),
			Qos:       mq.QosLevel(qos),
			Payload:   C.GoBytes(unsafe.Pointer(payload), payloadSize),
		})
	}
}

// Libmqtt_unsubscribe (int id, char *topic)
//export Libmqtt_unsubscribe
func Libmqtt_unsubscribe(id C.int, topic *C.char) {
	cid := int(id)
	if c, ok := clients[cid]; ok {
		c.UnSubscribe(C.GoString(topic))
	}
}

// Libmqtt_wait (int id)
//export Libmqtt_wait
func Libmqtt_wait(id C.int) {
	if c, ok := clients[int(id)]; ok {
		c.Wait()
	}
}

// Libmqtt_destroy (int id, bool force)
//export Libmqtt_destroy
func Libmqtt_destroy(id C.int, force C.bool) {
	if c, ok := clients[int(id)]; ok {
		c.Destroy(bool(force))
	}
}

// Libmqtt_set_pub_handler (int id, libmqtt_pub_handler h)
//export Libmqtt_set_pub_handler
func Libmqtt_set_pub_handler(id C.int, h C.libmqtt_pub_handler) {
	if c, ok := clients[int(id)]; ok {
		c.HandlePub(func(topic string, err error) {
			var er *C.char
			if err != nil {
				er = C.CString(err.Error())
			}
			C.call_pub_handler(h, C.CString(topic), er)
		})
	}
}

// Libmqtt_set_sub_handler (int id, libmqtt_sub_handler h)
//export Libmqtt_set_sub_handler
func Libmqtt_set_sub_handler(id C.int, h C.libmqtt_sub_handler) {
	if c, ok := clients[int(id)]; ok {
		c.HandleSub(func(topics []*mq.Topic, err error) {
			for _, t := range topics {
				var er *C.char
				if err != nil {
					er = C.CString(err.Error())
				}
				C.call_sub_handler(h, C.CString(t.Name), C.int(t.Qos), er)
			}
		})
	}
}

// Libmqtt_set_unsub_handler (int id, libmqtt_unsub_handler h)
//export Libmqtt_set_unsub_handler
func Libmqtt_set_unsub_handler(id C.int, h C.libmqtt_unsub_handler) {
	if c, ok := clients[int(id)]; ok {
		c.HandleUnSub(func(topics []string, err error) {
			for _, t := range topics {
				var er *C.char
				if err != nil {
					er = C.CString(err.Error())
				}
				C.call_unsub_handler(h, C.CString(t), er)
			}
		})
	}
}

// Libmqtt_set_net_handler (int id, libmqtt_net_handler h)
//export Libmqtt_set_net_handler
func Libmqtt_set_net_handler(id C.int, h C.libmqtt_net_handler) {
	if c, ok := clients[int(id)]; ok {
		c.HandleNet(func(server string, err error) {
			if err != nil {
				C.call_net_handler(h, C.CString(server), C.CString(err.Error()))
			}
		})
	}
}

// Libmqtt_set_persist_handler (int id, libmqtt_persist_handler h)
//export Libmqtt_set_persist_handler
func Libmqtt_set_persist_handler(id C.int, h C.libmqtt_persist_handler) {
	if c, ok := clients[int(id)]; ok {
		c.HandlePersist(func(err error) {
			if err != nil {
				C.call_persist_handler(h, C.CString(err.Error()))
			}
		})
	}
}

func wrapTopicHandler(h C.libmqtt_topic_handler) mq.TopicHandler {
	return func(topic string, qos mq.QosLevel, msg []byte) {
		C.call_topic_handler(h, C.CString(topic), C.int(qos), (*C.char)(C.CBytes(msg)), C.int(len(msg)))
	}
}

func main() {}
