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
	"bytes"
	"net/url"
	"sync"
	"sync/atomic"
	"testing"

	pah "github.com/eclipse/paho.mqtt.golang"
	lib "github.com/goiiot/libmqtt"
)

//smqM "github.com/surgemq/message"
//smq "github.com/surgemq/surgemq/service"
//gmq "github.com/yosssi/gmq/mqtt/client"

const (
	keepalive = 10
	server    = "localhost:1883"
	username  = "foo"
	password  = "bar"
	topic     = "/foo"
	qos       = 0

	pubCount = 50000
)

var (
	topicMsg = []byte("bar")
)

func BenchmarkLibmqttClient(b *testing.B) {
	b.N = pubCount
	b.ReportAllocs()
	var count uint32
	var client lib.Client
	var err error

	b.ResetTimer()
	if client, err = lib.NewClient(
		//lib.WithLog(lib.Verbose),
		lib.WithServer(server),
		lib.WithKeepalive(keepalive, 1.2),
		lib.WithIdentity(username, password),
		lib.WithRecvBuf(100),
		lib.WithSendBuf(100),
		lib.WithCleanSession(true)); err != nil {
		b.Log(err)
		b.FailNow()
	}

	client.HandleSub(func(topics []*lib.Topic, err error) {
		go func() {
			for i := 0; i < pubCount; i++ {
				client.Publish(&lib.PublishPacket{
					TopicName: topic,
					Qos:       qos,
					Payload:   topicMsg,
				})
			}
		}()
	})

	client.HandlePub(func(topic string, err error) {
		if err != nil {
			b.FailNow()
		}
	})

	client.HandleUnSub(func(topics []string, err error) {
		if err != nil {
			b.FailNow()
		}

		client.Destroy(true)
	})

	client.Handle(topic, func(t string, q lib.QosLevel, msg []byte) {
		if topic != t || q != qos || bytes.Compare(msg, topicMsg) != 0 {
			b.FailNow()
		}

		atomic.AddUint32(&count, 1)
		if atomic.LoadUint32(&count) == pubCount {
			client.UnSubscribe(topic)
		}
	})

	client.Connect(func(server string, code lib.ConnAckCode, err error) {
		if err != nil {
			b.Log(err)
			b.FailNow()
		} else if code != lib.ConnAccepted {
			b.FailNow()
		}

		client.Subscribe(&lib.Topic{Name: topic, Qos: qos})
	})

	client.Wait()
}

func BenchmarkPahoClient(b *testing.B) {
	b.N = pubCount
	b.ReportAllocs()
	var count uint32

	wg := &sync.WaitGroup{}
	wg.Add(1)
	b.ResetTimer()
	serverURL, err := url.Parse("tcp://" + server)
	if err != nil {
		b.FailNow()
	}

	client := pah.NewClient(&pah.ClientOptions{
		Servers:         []*url.URL{serverURL},
		KeepAlive:       keepalive,
		Username:        username,
		Password:        password,
		CleanSession:    true,
		ProtocolVersion: 4,
		Store:           pah.NewMemoryStore(),
	})

	t := client.Connect()
	if !t.Wait() {
		b.FailNow()
	}

	if err := t.Error(); err != nil {
		b.Log(err)
		b.FailNow()
	}

	t = client.Subscribe(topic, 0, func(c pah.Client, message pah.Message) {
		if topic != message.Topic() ||
			bytes.Compare(topicMsg, message.Payload()) != 0 ||
			qos != message.Qos() {
			b.FailNow()
		}

		atomic.AddUint32(&count, 1)
		if atomic.LoadUint32(&count) == pubCount {
			t := c.Unsubscribe(topic)
			if !t.Wait() {
				b.FailNow()
			}
			if err := t.Error(); err != nil {
				b.Log(err)
				b.FailNow()
			}

			c.Disconnect(0)
			wg.Done()
		}
	})

	if !t.Wait() {
		b.FailNow()
	}
	if err := t.Error(); err != nil {
		b.Log(err)
		b.FailNow()
	}

	for i := 0; i < pubCount; i++ {
		t = client.Publish(topic, 0, false, topicMsg)
		if !t.Wait() {
			b.FailNow()
		}

		if err := t.Error(); err != nil {
			b.Log(err)
			b.FailNow()
		}
	}
	wg.Wait()
}

//func BenchmarkGmqClient(b *testing.B) {
//	client := gmq.New(&gmq.Options{})
//	client.Connect(&gmq.ConnectOptions{
//		Network:   "tcp",
//		Address:   server,
//		UserName:  []byte(username),
//		Password:  []byte(password),
//		KeepAlive: keepalive,
//	})
//	subHandler := func(topicName, message []byte) {
//
//	}
//	if err := client.Subscribe(&gmq.SubscribeOptions{
//		SubReqs: []*gmq.SubReq{
//			{
//				TopicFilter: []byte(topic),
//				QoS:         qos,
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
//	err := client.Connect(server, &smqM.ConnectMessage{})
//	if err != nil {
//		b.Log(err)
//		b.FailNow()
//	}
//}
