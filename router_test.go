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

package libmqtt

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestTextRouter_Dispatch(t *testing.T) {
	r := &TextRouter{}
	count := 0
	testTopics := make(map[string]string)
	const topicCount = 1000

	addTopic := func() string {
		newTopic := strconv.Itoa(rand.New(rand.NewSource(time.Now().UnixNano())).Int())
		for _, ok := testTopics[newTopic]; ok; _, ok = testTopics[newTopic] {
			newTopic = strconv.Itoa(rand.New(rand.NewSource(time.Now().Unix())).Int())
		}
		testTopics[newTopic] = newTopic
		return newTopic
	}

	for i := 0; i < topicCount; i++ {
		newTopic := addTopic()
		r.Handle(newTopic, func(topic string, code SubAckCode, msg []byte) {
			if topic != newTopic {
				t.Log("fail at topic =", topic, ", target topic =", newTopic)
				t.FailNow()
			}
			count++
		})
	}

	for key := range testTopics {
		k := key
		r.Dispatch(&PublishPacket{TopicName: k})
	}

	if count != topicCount {
		t.Log("dispatch failed, count =", count)
		t.FailNow()
	}
}

func TestRegexRouter_Dispatch(t *testing.T) {
	r := &RegexRouter{}
	allCount, prefixCount, numCount := 0, 0, 0

	r.Handle(".*", func(topic string, code SubAckCode, msg []byte) {
		// should match all topics
		allCount++
	})

	r.Handle(`^(\/test)`, func(topic string, code SubAckCode, msg []byte) {
		// should match topics with `/test` prefix
		prefixCount++
	})

	r.Handle("\\d+", func(topic string, code SubAckCode, msg []byte) {
		// should match topics with number(s)
		numCount++
	})

	pkts := []*PublishPacket{
		{TopicName: "/test"},
		{TopicName: "/test/123"},
		{TopicName: "help"},
		{TopicName: "123"},
	}

	for _, v := range pkts {
		r.Dispatch(v)
	}

	if allCount != len(pkts) {
		t.Log("fail at all pkt count")
		t.FailNow()
	}

	if prefixCount != 2 {
		t.Log("fail at prefix pkt count")
		t.FailNow()
	}

	if numCount != 2 {
		t.Log("fail at num pkt count")
		t.FailNow()
	}
}

func TestRestRouter_Dispatch(t *testing.T) {

}
