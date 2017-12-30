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
