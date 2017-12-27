package main

// #cgo CFLAGS: -I include
/*
#include "libmqtth.h"

static inline void callConnHandler(libmqtt_conn_handler h,
	const char * server, libmqtt_connack_t code, const char * err) {
  if (h != NULL) {
    h(server, code, err);
  }
}

static inline void callPubHandler(libmqtt_pub_handler h,
	const char * topic, const char * msg) {
  if (h != NULL) {
    h(topic, msg);
  }
}

static inline void callSubHandler(libmqtt_sub_handler h,
	const char * topic, int qos, const char * err) {
  if (h != NULL) {
    h(topic, qos, err);
  }
}

static inline void callUnSubHandler(libmqtt_unsub_handler h,
	const char * topic, const char * err) {
  if (h != NULL) {
    h(topic, err);
  }
}

static inline void callNetHandler(libmqtt_net_handler h,
	const char * server, const char * err) {
  if (h != NULL) {
    h(server, err);
  }
}

static inline void callTopicHandler(libmqtt_topic_handler h,
	const char * topic, int qos , const char * msg) {
  if (h != NULL) {
    h(topic, qos, msg);
  }
}
*/
import "C"
import (
	"github.com/goiiot/libmqtt"
)

var (
	mConnHandler   C.libmqtt_conn_handler
	mPubHandler    C.libmqtt_pub_handler
	mSubHandler    C.libmqtt_sub_handler
	mUnSubHandler  C.libmqtt_unsub_handler
	mNetHandler    C.libmqtt_net_handler
	mTopicHandlers map[string]C.libmqtt_topic_handler
)

//export SetConnHandler
func SetConnHandler(h C.libmqtt_conn_handler) {
	mConnHandler = h
}

//export UnSetConnHandler
func UnSetConnHandler() {
	mConnHandler = nil
}

func callConnHandler(server string, code libmqtt.ConnAckCode, err error) {
	if mConnHandler != nil {
		var c C.libmqtt_connack_t
		switch code {
		case libmqtt.ConnAccepted:
			c = C.libmqtt_connack_accepted
		case libmqtt.ConnBadProtocol:
			c = C.libmqtt_connack_bad_proto
		case libmqtt.ConnIDRejected:
			c = C.libmqtt_connack_id_rejected
		case libmqtt.ConnServerUnavailable:
			c = C.libmqtt_connack_srv_unavail
		case libmqtt.ConnBadIdentity:
			c = C.libmqtt_connack_bad_identity
		case libmqtt.ConnAuthFail:
			c = C.libmqtt_connack_auth_fail
		}

		C.callConnHandler(mConnHandler, C.CString(server), c, C.CString(err.Error()))
	}
}

//export SetPubHandler
func SetPubHandler(h C.libmqtt_pub_handler) {
	mPubHandler = h
}

//export UnSetPubHandler
func UnSetPubHandler() {
	mPubHandler = nil
}

func callPubHandler(topic string, err error) {
	if mPubHandler != nil {
		C.callPubHandler(mPubHandler, C.CString(topic), C.CString(err.Error()))
	}
}

//export SetSubHandler
func SetSubHandler(h C.libmqtt_sub_handler) {
	mSubHandler = h
}

//export UnSubHandler
func UnSubHandler() {
	mSubHandler = nil
}

func callSubHandler(topics []*libmqtt.Topic, err error) {
	if mSubHandler != nil {
		for _, t := range topics {
			C.callSubHandler(mSubHandler, C.CString(t.Name), C.int(t.Qos), C.CString(err.Error()))
		}
	}
}

//export SetUnSubHandler
func SetUnSubHandler(h C.libmqtt_unsub_handler) {
	mUnSubHandler = h
}

//export UnSetUnSubHandler
func UnSetUnSubHandler() {
	mUnSubHandler = nil
}

func callUnSubHandler(topics []string, err error) {
	if mSubHandler != nil {
		for _, t := range topics {
			C.callUnSubHandler(mUnSubHandler, C.CString(t), C.CString(err.Error()))
		}
	}
}

//export SetNetHandler
func SetNetHandler(h C.libmqtt_net_handler) {
	mNetHandler = h
}

//export UnSetNetHandler
func UnSetNetHandler() {
	mNetHandler = nil
}

func callNetHandler(server string, err error) {
	if mNetHandler != nil {
		C.callNetHandler(mNetHandler, C.CString(server), C.CString(err.Error()))
	}
}

//export SetTopicHandler
func SetTopicHandler(topic *C.char, h C.libmqtt_topic_handler) {
	mTopicHandlers[C.GoString(topic)] = h
}

//export UnSetTopicHandler
func UnSetTopicHandler(topic *C.char) {

}

func callTopicHandler(topic string, qos libmqtt.QosLevel, msg []byte) {

}
