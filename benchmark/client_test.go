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
	gmq "github.com/yosssi/gmq/mqtt/client"
)

//smqM "github.com/surgemq/message"
//smq "github.com/surgemq/surgemq/service"

const (
	testKeepalive = 3600             // prevent keepalive packet disturb
	testServer    = "localhost:1883" // emqttd server address
	testUsername  = "foo"
	testPassword  = "bar"
	testTopic     = "/foo"
	testQos       = 0
	testBufSize   = 1024 // same with gmq default

	testPubCount = 10000
)

var (
	// 256 bytes
	testTopicMsg = []byte(
		"1234567890" + "1234567890" + "1234567890" + "1234567890" + "1234567890" +
			"1234567890" + "1234567890" + "1234567890" + "1234567890" + "1234567890" +
			"1234567890" + "1234567890" + "1234567890" + "1234567890" + "1234567890" +
			"1234567890" + "1234567890" + "1234567890" + "1234567890" + "1234567890" +
			"1234567890" + "1234567890" + "1234567890" + "1234567890" + "1234567890" +
			"12345")
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
		lib.WithRecvBuf(testBufSize),
		lib.WithSendBuf(testBufSize),
		lib.WithCleanSession(true)); err != nil {
		b.Log(err)
		b.FailNow()
	}

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

		for i := 0; i < b.N; i++ {
			client.Publish(&lib.PublishPacket{
				TopicName: testTopic,
				Qos:       testQos,
				Payload:   testTopicMsg,
			})
		}
	})
	b.ResetTimer()
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
		Servers:             []*url.URL{serverURL},
		KeepAlive:           testKeepalive,
		Username:            testUsername,
		Password:            testPassword,
		CleanSession:        true,
		ProtocolVersion:     4,
		MessageChannelDepth: testBufSize,
		Store:               pah.NewMemoryStore(),
	})

	t := client.Connect()
	if !t.Wait() {
		b.FailNow()
	}

	if err := t.Error(); err != nil {
		b.Log(err)
		b.FailNow()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.Publish(testTopic, 0, false, testTopicMsg)
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

func BenchmarkGmqClient(b *testing.B) {
	b.N = testPubCount
	b.ReportAllocs()

	client := gmq.New(&gmq.Options{ErrorHandler: func(e error) {
		if e != nil {
			b.Log(e)
			b.FailNow()
		}
	}})
	if err := client.Connect(&gmq.ConnectOptions{
		Network:      "tcp",
		Address:      testServer,
		UserName:     []byte(testUsername),
		Password:     []byte(testPassword),
		KeepAlive:    testKeepalive,
		CleanSession: true,
	}); err != nil {
		b.Log(err)
		b.FailNow()
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := client.Publish(&gmq.PublishOptions{
			QoS:       testQos,
			Retain:    false,
			TopicName: []byte(testTopic),
			Message:   testTopicMsg,
		}); err != nil {
			b.Log(err)
			b.FailNow()
		}
	}
	if err := client.Unsubscribe(&gmq.UnsubscribeOptions{
		TopicFilters: [][]byte{[]byte(testTopic)},
	}); err != nil {
		b.Log(err)
		b.FailNow()
	}

	client.Terminate()
}

//func BenchmarkSurgeClient(b *testing.B) {
//	client := &smq.Client{KeepAlive: 10}
//	err := client.Connect(testServer, &smqM.ConnectMessage{})
//	if err != nil {
//		b.Log(err)
//		b.FailNow()
//	}
//}
