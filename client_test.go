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
)

// test with emqttd server (http://emqtt.io/ or https://github.com/emqtt/emqttd)
// the server is configured with default configuration

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

type extraHandler struct {
	pubH   func()
	subH   func()
	unSubH func()
	netH   func()
}

func plainClient(t *testing.T, exH *extraHandler) Client {
	c, err := NewClient(
		WithServer("localhost:1883"),
		WithDialTimeout(10),
		WithKeepalive(10, 1.2),
		WithIdentity("admin", "public"),
		WithWill("test", Qos0, false, []byte("test data")),
		WithLog(Verbose),
	)

	if err != nil {
		t.Log(err)
		t.Log("create client failed")
		t.FailNow()
	}
	initClient(c, exH, t)
	return c
}

func tlsClient(t *testing.T, exH *extraHandler) Client {
	c, err := NewClient(
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
		WithLog(Verbose),
	)

	if err != nil {
		t.Log(err)
		t.Log("create client failed")
		t.FailNow()
	}
	initClient(c, exH, t)
	return c
}

func initClient(c Client, exH *extraHandler, t *testing.T) {
	var pubCount uint32
	c.HandlePub(func(topic string, err error) {
		if err != nil {
			t.Log("pub failed, err =", err)
			t.FailNow()
		}
		atomic.AddUint32(&pubCount, 1)
		if atomic.LoadUint32(&pubCount) == uint32(len(testMsgs)) {
			if exH != nil && exH.pubH != nil {
				exH.pubH()
			}
		}
	})

	c.HandleSub(func(topics []*Topic, err error) {
		if err != nil {
			t.Log("sub failed, err =", err)
			t.FailNow()
		}

		if exH != nil && exH.subH != nil {
			exH.subH()
		}
	})

	c.HandleUnSub(func(topics []string, err error) {
		if err != nil {
			t.Log("unsub failed, err =", err)
			t.FailNow()
		}

		if exH != nil && exH.unSubH != nil {
			exH.unSubH()
		}
	})

	c.HandleNet(func(server string, err error) {
		if err != nil {
			t.Log("net error, err =", err)
			t.FailNow()
		}

		if exH != nil && exH.netH != nil {
			exH.netH()
		}
	})
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

func testSub(c Client, t *testing.T) {
	var count uint32
	N := uint32(len(testMsgs))
	for index, value := range testMsgs {
		i := uint32(index)
		v := value
		c.Handle(v.TopicName, func(topic string, maxQos SubAckCode, msg []byte) {
			if maxQos != v.Qos || bytes.Compare(v.Payload, msg) != 0 {
				t.Log("fail at sub topic =", topic,
					", content unexcepted, payload =", string(msg),
					"target payload =", string(v.Payload),
				)
				t.FailNow()
			} else {
				atomic.AddUint32(&count, 1)
			}

			if i == N-1 {
				if atomic.LoadUint32(&count) != N {
					t.Log("sub recv count =", c)
					t.FailNow()
				}
			}
		})
	}
	c.Subscribe(testSubs...)
}

func TestNewClient(t *testing.T) {
	_, err := NewClient()
	if err == nil {
		t.Log("empty server not failed")
		t.FailNow()
	}

	_, err = NewClient(
		WithTLS(
			"foo",
			"bar",
			"foobar",
			"foo.bar",
			true),
	)
	if err == nil {
		t.Log("tls not failed")
		t.FailNow()
	}
}

// conn
func TestClient_Connect(t *testing.T) {
	var c Client
	afterConn := func() {
		c.Destroy(true)
	}

	c = plainClient(t, nil)
	testConn(c, t, afterConn)
	c.Wait()

	c = tlsClient(t, nil)
	testConn(c, t, afterConn)
	c.Wait()
}

// conn -> pub
func TestClient_Publish(t *testing.T) {
	var c Client
	afterConn := func() {
		c.Publish(testMsgs...)
	}

	exH := &extraHandler{
		pubH: func() {
			c.Destroy(true)
		},
	}

	c = plainClient(t, exH)
	testConn(c, t, afterConn)
	c.Wait()

	c = tlsClient(t, exH)
	testConn(c, t, afterConn)
	c.Wait()
}

// conn -> sub -> pub
func TestClient_Subscribe(t *testing.T) {
	var c Client
	afterConn := func() {
		testSub(c, t)
	}

	extH := &extraHandler{
		subH: func() {
			c.Publish(testMsgs...)
		},
		pubH: func() {
			c.Destroy(true)
		},
	}

	c = plainClient(t, extH)
	testConn(c, t, afterConn)
	c.Wait()

	c = tlsClient(t, extH)
	testConn(c, t, afterConn)
	c.Wait()
}

// conn -> sub -> pub -> unSub
func TestClient_UnSubscribe(t *testing.T) {
	var c Client
	afterConn := func() {
		testSub(c, t)
	}

	extH := &extraHandler{
		subH: func() {
			c.UnSubscribe(testTopics...)
		},
		unSubH: func() {
			c.Destroy(true)
		},
	}

	c = plainClient(t, extH)
	testConn(c, t, afterConn)
	c.Wait()

	c = tlsClient(t, extH)
	testConn(c, t, afterConn)
	c.Wait()
}

func TestClient_Reconnect(t *testing.T) {
	c := plainClient(t, nil).(*client)
	count := &atomic.Value{}
	count.Store(0)
	c.Connect(func(server string, code ConnAckCode, err error) {
		if err != nil {
			t.Log(err)
			t.FailNow()
		}

		if code != ConnAccepted {
			t.Log(code)
			t.FailNow()
		}

		if count.Load().(int) < 3 {
			c.conn.Range(func(key, value interface{}) bool {
				v := value.(*connImpl)
				v.conn.Close()
				return true
			})
			count.Store(count.Load().(int) + 1)
		} else {
			c.Destroy(true)
		}
	})
	c.Wait()
}
