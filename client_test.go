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
	"bytes"
	"sync/atomic"
	"testing"
	"time"
)

// test with emqttd server (http://emqtt.io/ or https://github.com/emqtt/emqttd)
// the server is configured with default configuration

const (
	waitTime = 1 * time.Second
)

var (
	testTopics = []string{"/test", "/test/foo", "/test/bar"}
	testMsgs   = []*PublishPacket{
		{TopicName: testTopics[0], Qos: Qos0, Payload: []byte("test")},
		{TopicName: testTopics[1], Qos: Qos1, Payload: []byte("foo")},
		{TopicName: testTopics[2], Qos: Qos2, Payload: []byte("bar")},
	}
	testSubs = []*Topic{
		{Name: testTopics[0], Qos: Qos0},
		{Name: testTopics[1], Qos: Qos1},
		{Name: testTopics[2], Qos: Qos2},
	}
)

func plainClient() Client {
	return NewClient(
		WithServer("localhost:1883"),
		WithDialTimeout(10),
		WithKeepalive(10, 1.2),
		WithIdentity("admin", "public"),
		WithWill("test", Qos0, false, []byte("test data")),
		WithLogger(Verbose),
	)
}

func tlsClient() Client {
	return NewClient(
		WithServer("localhost:8883"),
		WithTLS(
			"./testdata/client-cert.pem",
			"./testdata/client-key.pem",
			"./testdata/ca-cert.pem",
			"MacBook-Air.local",
			true),
		WithDialTimeout(10),
		WithKeepalive(10, 1.2),
		WithIdentity("admin", "public"),
		WithWill("test", Qos0, false, []byte("test data")),
		WithLogger(Verbose),
	)
}

// conn
func TestClient_Connect(t *testing.T) {
	var c Client
	afterConn := func() {
		c.Destroy(true)
	}

	c = plainClient()
	testConn(c, t, afterConn)
	c.Wait()

	c = tlsClient()
	testConn(c, t, afterConn)
	c.Wait()
}

func testConn(c Client, t *testing.T, afterConn func()) {
	c.Connect(func(server string, code ConnAckCode, err error) {
		if err != nil {
			t.Log(err)
			t.FailNow()
		}

		if code != ConnAccepted {
			t.Log(code)
			t.FailNow()
		}

		if afterConn != nil {
			afterConn()
		}
	})
}

// conn -> pub
func TestClient_Publish(t *testing.T) {
	var c Client
	afterConn := func() {
		c.Publish(testMsgs...)
		<-time.After(waitTime)
		c.Destroy(true)
	}

	c = plainClient()
	testConn(c, t, afterConn)
	c.Wait()

	c = tlsClient()
	testConn(c, t, afterConn)
	c.Wait()
}

func TestClient_Subscribe(t *testing.T) {
	var c Client
	afterConn := func() {
		testSub(c, t)
		<-time.After(waitTime)
		c.Publish(testMsgs...)
		<-time.After(waitTime)
		c.Destroy(true)
	}

	c = plainClient()
	testConn(c, t, afterConn)
	c.Wait()

	c = tlsClient()
	testConn(c, t, afterConn)
	c.Wait()
}

func testSub(c Client, t *testing.T) {
	count := atomic.Value{}
	count.Store(int(0))
	N := len(testMsgs)
	for index, value := range testMsgs {
		i := index
		v := value
		c.Handle(v.TopicName, func(topic string, maxQos SubAckCode, msg []byte) {
			if maxQos != v.Qos || bytes.Compare(v.Payload, msg) != 0 {
				t.Log("fail at sub topic =", topic,
					", content unexcepted, payload =", string(msg),
					"target payload =", string(v.Payload),
				)
				t.FailNow()
			} else {
				count.Store(count.Load().(int) + 1)
			}

			if i == N-1 {
				if c := count.Load().(int); c != N {
					t.Log("sub recv count =", c)
					t.FailNow()
				}
			}
		})
	}
	c.Subscribe(testSubs...)
}

// conn -> sub -> pub -> unSub
func TestClient_UnSubscribe(t *testing.T) {
	var c Client
	afterConn := func() {
		testSub(c, t)
		<-time.After(waitTime)
		c.Publish(testMsgs...)
		<-time.After(waitTime)
		c.UnSubscribe(testTopics...)
		<-time.After(waitTime)
		c.Publish(testMsgs...)
		<-time.After(waitTime)
		c.Destroy(true)
	}

	c = plainClient()
	testConn(c, t, afterConn)
	c.Wait()

	c = tlsClient()
	testConn(c, t, afterConn)
	c.Wait()
}
