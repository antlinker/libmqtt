//+build cgo lib

package main

// #cgo CFLAGS: -I include
/*
#include "LibMQTTConst.h"
#include "LibMQTTCallback.h"
*/
import "C"

//export Publish
// Publish(topic *C.char, qos C.goiiot_lib_mqtt_qos_t, payload *C.char, payloadSize C.int)
func Publish(topic *C.char, qos C.goiiot_lib_mqtt_qos_t, payload *C.char, payloadSize C.int) {
	// goTopic := C.GoString(topic)
}

func main() {}
