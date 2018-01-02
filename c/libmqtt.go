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
#include <stdlib.h>
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
  (int client, char *server, libmqtt_connack_t code, char *err);

typedef void (*libmqtt_pub_handler)
  (int client, char *topic, char *err);

typedef void (*libmqtt_sub_handler)
  (int client, char *topic, int qos, char *err);

typedef void (*libmqtt_unsub_handler)
  (int client, char *topic, char *err);

typedef void (*libmqtt_net_handler)
  (int client, char *server, char *err);

typedef void (*libmqtt_persist_handler)
  (int client, char *err);

typedef void (*libmqtt_topic_handler)
  (int client, char *topic, int qos, char *msg, int size);

static inline void call_conn_handler
  (libmqtt_conn_handler h, int client, char * server, libmqtt_connack_t code, char * err) {
  if (h != NULL) {
    h(client, server, code, err);
    free(server);
    free(err);
  }
}

static inline void call_pub_handler
  (libmqtt_pub_handler h, int client, char * topic, char * err) {
  if (h != NULL) {
    h(client, topic, err);
    free(topic);
    free(err);
  }
}

static inline void call_sub_handler
  (libmqtt_sub_handler h, int client, char * topic, int qos, char * err) {
  if (h != NULL) {
    h(client, topic, qos, err);
    free(topic);
    free(err);
  }
}

static inline void call_unsub_handler
  (libmqtt_unsub_handler h, int client, char * topic, char * err) {
  if (h != NULL) {
    h(client, topic, err);
    free(topic);
    free(err);
  }
}

static inline void call_net_handler
  (libmqtt_net_handler h, int client, char * server, char * err) {
  if (h != NULL) {
    h(client, server, err);
    free(server);
    free(err);
  }
}

static inline void call_persist_handler
  (libmqtt_persist_handler h, int client, char *err) {
  if (h != NULL) {
    h(client, err);
    free(err);
  }
}

static inline void call_topic_handler
  (libmqtt_topic_handler h, int client, char * topic, int qos , char * msg, int size) {
  if (h != NULL) {
    h(client, topic, qos, msg, size);
    free(topic);
    free(msg);
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

// Libmqtt_client_set_server (int client, char *server)
//export Libmqtt_client_set_server
func Libmqtt_client_set_server(client C.int, server *C.char) {
	addOption(client, mq.WithServer(C.GoString(server)))
}

// Libmqtt_client_set_none_persist (int client)
//export Libmqtt_client_set_none_persist
func Libmqtt_client_set_none_persist(client C.int) {
	addOption(client, mq.WithPersist(&mq.NonePersist{}))
}

// Libmqtt_client_set_mem_persist (int client, int maxCount, bool dropOnExceed, bool duplicateReplace)
//export Libmqtt_client_set_mem_persist
func Libmqtt_client_set_mem_persist(client C.int, maxCount C.int, dropOnExceed C.bool, duplicateReplace C.bool) {
	addOption(client, mq.WithPersist(mq.NewMemPersist(&mq.PersistStrategy{
		MaxCount:         uint32(maxCount),
		DropOnExceed:     bool(dropOnExceed),
		DuplicateReplace: bool(duplicateReplace),
	})))
}

// Libmqtt_client_set_file_persist (int client, char *dirPath, int maxCount, bool dropOnExceed, bool duplicateReplace)
//export Libmqtt_client_set_file_persist
func Libmqtt_client_set_file_persist(client C.int, dirPath *C.char, maxCount C.int, dropOnExceed C.bool, duplicateReplace C.bool) {
	addOption(client, mq.WithPersist(mq.NewFilePersist(C.GoString(dirPath), &mq.PersistStrategy{
		MaxCount:         uint32(maxCount),
		DropOnExceed:     bool(dropOnExceed),
		DuplicateReplace: bool(duplicateReplace),
	})))
}

// Libmqtt_client_set_clean_session (bool flag)
//export Libmqtt_client_set_clean_session
func Libmqtt_client_set_clean_session(client C.int, flag C.bool) {
	addOption(client, mq.WithCleanSession(bool(flag)))
}

// Libmqtt_client_set_client_id (int client, char * client_id)
//export Libmqtt_client_set_client_id
func Libmqtt_client_set_client_id(client C.int, clientID *C.char) {
	addOption(client, mq.WithClientID(C.GoString(clientID)))
}

// Libmqtt_client_set_dial_timeout (int client, int timeout)
//export Libmqtt_client_set_dial_timeout
func Libmqtt_client_set_dial_timeout(client C.int, timeout C.int) {
	addOption(client, mq.WithDialTimeout(uint16(timeout)))
}

// Libmqtt_client_set_identity (int client, char * username, char * password)
//export Libmqtt_client_set_identity
func Libmqtt_client_set_identity(client C.int, username, password *C.char) {
	addOption(client, mq.WithIdentity(C.GoString(username), C.GoString(password)))
}

// Libmqtt_client_set_keepalive (int client, int keepalive, float factor)
//export Libmqtt_client_set_keepalive
func Libmqtt_client_set_keepalive(client C.int, keepalive C.int, factor C.float) {
	addOption(client, mq.WithKeepalive(uint16(keepalive), float64(factor)))
}

// Libmqtt_client_set_log (int client, libmqtt_log_level l)
//export Libmqtt_client_set_log
func Libmqtt_client_set_log(client C.int, l C.libmqtt_log_level) {
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
	addOption(client, mq.WithLog(level))
}

// Libmqtt_client_set_recv_buf (int client, int size)
//export Libmqtt_client_set_recv_buf
func Libmqtt_client_set_recv_buf(client C.int, size C.int) {
	addOption(client, mq.WithRecvBuf(int(size)))
}

// Libmqtt_client_set_send_buf (int client, int size)
//export Libmqtt_client_set_send_buf
func Libmqtt_client_set_send_buf(client C.int, size C.int) {
	addOption(client, mq.WithSendBuf(int(size)))
}

// Libmqtt_client_set_tls (int client, char * certFile, char * keyFile, char * caCert, char * serverNameOverride, bool skipVerify)
// use ssl to connect
//export Libmqtt_client_set_tls
func Libmqtt_client_set_tls(client C.int, certFile, keyFile, caCert, serverNameOverride *C.char, skipVerify C.bool) {
	addOption(client, mq.WithTLS(
		C.GoString(certFile),
		C.GoString(keyFile),
		C.GoString(caCert),
		C.GoString(serverNameOverride),
		bool(skipVerify)))
}

// Libmqtt_client_set_will (int client, char *topic, int qos, bool retain, char *payload, int payloadSize)
// mark this connection with will message
//export Libmqtt_client_set_will
func Libmqtt_client_set_will(client C.int, topic *C.char, qos C.int, retain C.bool, payload *C.char, payloadSize C.int) {
	addOption(client, mq.WithWill(
		C.GoString(topic),
		mq.QosLevel(qos),
		bool(retain),
		C.GoBytes(unsafe.Pointer(payload), payloadSize)))
}

// Libmqtt_setup (int client)
// setup the client with previously defined options
//export Libmqtt_setup
func Libmqtt_setup(client C.int) *C.char {
	cid := int(client)
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

// Libmqtt_handle (int client, char *topic, libmqtt_topic_handler h)
//export Libmqtt_handle
func Libmqtt_handle(client C.int, topic *C.char, h C.libmqtt_topic_handler) {
	if c, ok := clients[int(client)]; ok {
		c.Handle(C.GoString(topic), wrapTopicHandler(client, h))
	}
}

// Libmqtt_connect (int client)
//export Libmqtt_connect
func Libmqtt_connect(client C.int, h C.libmqtt_conn_handler) {
	if c, ok := clients[int(client)]; ok {
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
			C.call_conn_handler(h, client, C.CString(server), c, er)
		})
	}
}

// Libmqtt_subscribe (int client, char *topic, int qos)
//export Libmqtt_subscribe
func Libmqtt_subscribe(client C.int, topic *C.char, qos C.int) {
	if c, ok := clients[int(client)]; ok {
		c.Subscribe(&mq.Topic{
			Name: C.GoString(topic),
			Qos:  mq.QosLevel(qos),
		})
	}
}

// Libmqtt_publish (int client, char *topic, int qos, char *payload, int payloadSize)
//export Libmqtt_publish
func Libmqtt_publish(client C.int, topic *C.char, qos C.int, payload *C.char, payloadSize C.int) {
	if c, ok := clients[int(client)]; ok {
		c.Publish(&mq.PublishPacket{
			TopicName: C.GoString(topic),
			Qos:       mq.QosLevel(qos),
			Payload:   C.GoBytes(unsafe.Pointer(payload), payloadSize),
		})
	}
}

// Libmqtt_unsubscribe (int client, char *topic)
//export Libmqtt_unsubscribe
func Libmqtt_unsubscribe(client C.int, topic *C.char) {
	cid := int(client)
	if c, ok := clients[cid]; ok {
		c.UnSubscribe(C.GoString(topic))
	}
}

// Libmqtt_wait (int client)
//export Libmqtt_wait
func Libmqtt_wait(client C.int) {
	if c, ok := clients[int(client)]; ok {
		c.Wait()
	}
}

// Libmqtt_destroy (int client, bool force)
//export Libmqtt_destroy
func Libmqtt_destroy(client C.int, force C.bool) {
	if c, ok := clients[int(client)]; ok {
		c.Destroy(bool(force))
	}
}

// Libmqtt_set_pub_handler (int client, libmqtt_pub_handler h)
//export Libmqtt_set_pub_handler
func Libmqtt_set_pub_handler(client C.int, h C.libmqtt_pub_handler) {
	if c, ok := clients[int(client)]; ok {
		c.HandlePub(func(topic string, err error) {
			var er *C.char
			if err != nil {
				er = C.CString(err.Error())
			}
			C.call_pub_handler(h, client, C.CString(topic), er)
		})
	}
}

// Libmqtt_set_sub_handler (int client, libmqtt_sub_handler h)
//export Libmqtt_set_sub_handler
func Libmqtt_set_sub_handler(client C.int, h C.libmqtt_sub_handler) {
	if c, ok := clients[int(client)]; ok {
		c.HandleSub(func(topics []*mq.Topic, err error) {
			for _, t := range topics {
				var er *C.char
				if err != nil {
					er = C.CString(err.Error())
				}
				C.call_sub_handler(h, client, C.CString(t.Name), C.int(t.Qos), er)
			}
		})
	}
}

// Libmqtt_set_unsub_handler (int client, libmqtt_unsub_handler h)
//export Libmqtt_set_unsub_handler
func Libmqtt_set_unsub_handler(client C.int, h C.libmqtt_unsub_handler) {
	if c, ok := clients[int(client)]; ok {
		c.HandleUnSub(func(topics []string, err error) {
			for _, t := range topics {
				var er *C.char
				if err != nil {
					er = C.CString(err.Error())
				}
				C.call_unsub_handler(h, client, C.CString(t), er)
			}
		})
	}
}

// Libmqtt_set_net_handler (int client, libmqtt_net_handler h)
//export Libmqtt_set_net_handler
func Libmqtt_set_net_handler(client C.int, h C.libmqtt_net_handler) {
	if c, ok := clients[int(client)]; ok {
		c.HandleNet(func(server string, err error) {
			if err != nil {
				C.call_net_handler(h, client, C.CString(server), C.CString(err.Error()))
			}
		})
	}
}

// Libmqtt_set_persist_handler (int client, libmqtt_persist_handler h)
//export Libmqtt_set_persist_handler
func Libmqtt_set_persist_handler(client C.int, h C.libmqtt_persist_handler) {
	if c, ok := clients[int(client)]; ok {
		c.HandlePersist(func(err error) {
			if err != nil {
				C.call_persist_handler(h, client, C.CString(err.Error()))
			}
		})
	}
}

func wrapTopicHandler(client C.int, h C.libmqtt_topic_handler) mq.TopicHandler {
	return func(topic string, qos mq.QosLevel, msg []byte) {
		C.call_topic_handler(h, client, C.CString(topic), C.int(qos), (*C.char)(C.CBytes(msg)), C.int(len(msg)))
	}
}

func addOption(client C.int, option mq.Option) {
	cid := int(client)
	if v, ok := clientOptions[cid]; ok {
		v = append(v, option)
		clientOptions[cid] = v
	}
}

func main() {}
