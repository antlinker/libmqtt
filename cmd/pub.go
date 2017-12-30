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

func execPub(args []string) bool {
	if client == nil {
		println("please connect to server first")
		return true
	}

	pubs := make([]*mq.PublishPacket, 0)
	for _, v := range args {
		pubStr := strings.Split(v, ",")
		if len(pubStr) != 3 {
			pubUsage()
			return true
		}
		qos, err := strconv.Atoi(pubStr[1])
		if err != nil {
			pubUsage()
			return false
		}
		pubs = append(pubs, &mq.PublishPacket{
			TopicName: pubStr[0],
			Qos:       mq.QosLevel(qos),
			Payload:   []byte(pubStr[2]),
		})
	}
	client.Publish(pubs...)
	return true
}

func pubUsage() {
	println(`p, pub [topic,qos,message] [...] - publish topic message(s)`)
}
