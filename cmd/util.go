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

import (
	mq "github.com/goiiot/libmqtt"
)

func connHandler(server string, code mq.ConnAckCode, err error) {
	if err != nil {
		println("\nconnect to server error:", err)
	} else if code != mq.ConnAccepted {
		println("\nconnection rejected by server, code:", code)
	} else {
		println("\nconnected to server")
	}
	print(lineStart)
}

func pubHandler(topic string, err error) {
	if err != nil {
		println("\npub", topic, "failed, error =", err)
	} else {
		println("\npub", topic, "success")
	}
	print(lineStart)
}

func subHandler(topics []*mq.Topic, err error) {
	if err != nil {
		println("\nsub", topics, "failed, error =", err)
	} else {
		println("\nsub", topics, "success")
	}
	print(lineStart)
}

func unSubHandler(topics []string, err error) {
	if err != nil {
		println("\nunsub", topics, "failed, error =", err)
	} else {
		println("\nunsub", topics, "success")
	}
	print(lineStart)
}

func netHandler(server string, err error) {
	println("\nconnection to server, error:", err)
	print(lineStart)
}

func topicHandler(topic string, qos mq.QosLevel, msg []byte) {
	println("\n[MSG] topic:", topic, "msg:", string(msg), "qos:", qos)
	print(lineStart)
}

func invalidQos() {
	println("\nqos level should either be 0, 1 or 2")
	print(lineStart)
}
