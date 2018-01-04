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

package benchmark

import (
	"net/url"
	"sync/atomic"
	"testing"

	pah "github.com/eclipse/paho.mqtt.golang"
	lib "github.com/goiiot/libmqtt"
)

//smqM "github.com/surgemq/message"
//smq "github.com/surgemq/surgemq/service"
//gmq "github.com/yosssi/gmq/mqtt/client"

const (
	testKeepalive = 3600
	testServer    = "localhost:1883"
	testUsername  = "foo"
	testPassword  = "bar"
	testTopic     = "/foo"
	testQos       = 0

	testPubCount = 10000000
)

var (
	testTopicMsg = []byte("bar")
)

func BenchmarkLibmqttClient(b *testing.B) {
	b.N = testPubCount
	b.ReportAllocs()
	var count uint32
	var client lib.Client
	var err error

	if client, err = lib.NewClient(
		lib.WithServer(testServer),
		lib.WithKeepalive(testKeepalive, 1.2),
		lib.WithIdentity(testUsername, testPassword),
		lib.WithRecvBuf(100),
		lib.WithSendBuf(100),
		lib.WithCleanSession(true)); err != nil {
		b.Log(err)
		b.FailNow()
	}

	b.ResetTimer()
	client.HandleUnSub(func(topic []string, err error) {
		if err != nil {
			b.Log(err)
			b.FailNow()
		}
		client.Destroy(true)
	})
	client.HandlePub(func(topic string, err error) {
		if err != nil {
			b.FailNow()
		}
		atomic.AddUint32(&count, 1)
		if atomic.LoadUint32(&count) == testPubCount {
			client.UnSubscribe(topic)
		}
	})
	client.Connect(func(server string, code lib.ConnAckCode, err error) {
		if err != nil {
			b.Log(err)
			b.FailNow()
		} else if code != lib.ConnAccepted {
			b.Log(code)
			b.FailNow()
		}

		pubs := make([]*lib.PublishPacket, b.N)
		pkt := &lib.PublishPacket{
			TopicName: testTopic,
			Qos:       testQos,
			Payload:   testTopicMsg,
		}
		for i := 0; i < b.N; i++ {
			pubs[i] = pkt
		}
		client.Publish(pubs...)

	})
	client.Wait()
}

func BenchmarkPahoClient(b *testing.B) {
	b.N = testPubCount
	b.ReportAllocs()

	serverURL, err := url.Parse("tcp://" + testServer)
	if err != nil {
		b.Log(err)
		b.FailNow()
	}

	client := pah.NewClient(&pah.ClientOptions{
		Servers:         []*url.URL{serverURL},
		KeepAlive:       testKeepalive,
		Username:        testUsername,
		Password:        testPassword,
		CleanSession:    true,
		ProtocolVersion: 4,
		Store:           pah.NewMemoryStore(),
	})

	b.ResetTimer()
	t := client.Connect()
	if !t.Wait() {
		b.FailNow()
	}

	if err := t.Error(); err != nil {
		b.Log(err)
		b.FailNow()
	}

	if !t.Wait() {
		b.FailNow()
	}
	if err := t.Error(); err != nil {
		b.Log(err)
		b.FailNow()
	}

	for i := 0; i < b.N; i++ {
		t = client.Publish(testTopic, 0, false, testTopicMsg)
		if !t.Wait() {
			b.FailNow()
		}

		if err := t.Error(); err != nil {
			b.Log(err)
			b.FailNow()
		}
	}
	t = client.Unsubscribe(testTopic)
	if !t.Wait() {
		b.FailNow()
	}
	if err := t.Error(); err != nil {
		b.Log(err)
		b.FailNow()
	}

	client.Disconnect(0)
}

//func BenchmarkGmqClient(b *testing.B) {
//	client := gmq.New(&gmq.Options{})
//	client.Connect(&gmq.ConnectOptions{
//		Network:   "tcp",
//		Address:   testServer,
//		UserName:  []byte(testUsername),
//		Password:  []byte(testPassword),
//		KeepAlive: testKeepalive,
//	})
//	subHandler := func(topicName, message []byte) {
//
//	}
//	if err := client.Subscribe(&gmq.SubscribeOptions{
//		SubReqs: []*gmq.SubReq{
//			{
//				TopicFilter: []byte(testTopic),
//				QoS:         testQos,
//				Handler:     subHandler,
//			},
//		},
//	}); err != nil {
//		b.Log(err)
//		b.FailNow()
//	}
//
//}

//func BenchmarkSurgeClient(b *testing.B) {
//	client := &smq.Client{KeepAlive: 10}
//	err := client.Connect(testServer, &smqM.ConnectMessage{})
//	if err != nil {
//		b.Log(err)
//		b.FailNow()
//	}
//}
