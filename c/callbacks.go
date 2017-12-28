package main

// #cgo CFLAGS: -I include
/*
#include "libmqtth.h"
*/
import "C"

var (
	mConnHandler  C.libmqtt_conn_handler
	mPubHandler   C.libmqtt_pub_handler
	mSubHandler   C.libmqtt_sub_handler
	mUnSubHandler C.libmqtt_unsub_handler
	mNetHandler   C.libmqtt_net_handler
)

// SetConnHandler (libmqtt_conn_handler h)
//export SetConnHandler
func SetConnHandler(h C.libmqtt_conn_handler) {
	mConnHandler = h
}

// UnSetConnHandler ()
//export UnSetConnHandler
func UnSetConnHandler() {
	mConnHandler = nil
}

// SetPubHandler (libmqtt_pub_handler h)
//export SetPubHandler
func SetPubHandler(h C.libmqtt_pub_handler) {
	mPubHandler = h
}

// UnSetPubHandler ()
//export UnSetPubHandler
func UnSetPubHandler() {
	mPubHandler = nil
}

// SetSubHandler (libmqtt_sub_handler h)
//export SetSubHandler
func SetSubHandler(h C.libmqtt_sub_handler) {
	mSubHandler = h
}

// UnSubHandler ()
//export UnSubHandler
func UnSubHandler() {
	mSubHandler = nil
}

// SetUnSubHandler (libmqtt_unsub_handler h)
//export SetUnSubHandler
func SetUnSubHandler(h C.libmqtt_unsub_handler) {
	mUnSubHandler = h
}

// UnSetUnSubHandler ()
//export UnSetUnSubHandler
func UnSetUnSubHandler() {
	mUnSubHandler = nil
}

// SetNetHandler (libmqtt_net_handler h)
//export SetNetHandler
func SetNetHandler(h C.libmqtt_net_handler) {
	mNetHandler = h
}

// UnSetNetHandler ()
//export UnSetNetHandler
func UnSetNetHandler() {
	mNetHandler = nil
}
