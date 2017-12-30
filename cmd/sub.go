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
	"strconv"
	"strings"

	mq "github.com/goiiot/libmqtt"
)

func execSub(args []string) bool {
	if client == nil {
		println("please connect to server first")
		return true
	}

	topics := make([]*mq.Topic, 0)
	for _, v := range args {
		topicStr := strings.Split(v, ",")
		if len(topicStr) != 2 {
			subUsage()
			return true
		}
		qos, err := strconv.Atoi(topicStr[1])
		if err != nil {
			subUsage()
			return true
		}
		topics = append(topics, &mq.Topic{Name: topicStr[0], Qos: mq.QosLevel(qos)})
	}
	for _, t := range topics {
		client.Handle(t.Name, topicHandler)
	}
	client.Subscribe(topics...)
	return true
}

func execUnSub(args []string) bool {
	if client == nil {
		println("please connect to server first")
		return true
	}

	client.UnSubscribe(args...)
	return true
}

func subUsage() {
	println(`s, sub [topic,qos] [...] - subscribe topic(s)`)
}

func unSubUsage() {
	println(`u, unsub [topic] [...] - unsubscribe topic(s)`)
}
