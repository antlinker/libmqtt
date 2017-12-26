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
	waitTime = 2 * time.Second
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

func connHandler(t *testing.T) ConnHandler {
	return func(server string, code ConAckCode, err error) {
		if err != nil {
			t.Log(err)
			t.FailNow()
		}

		if code != ConnAccepted {
			t.Log(code)
			t.FailNow()
		}
	}
}

// conn
func TestClient_Connect(t *testing.T) {
	c := plainClient()
	testConn(c, t)
	<-time.After(waitTime)
	c.Destroy(true)

	c = tlsClient()
	testConn(c, t)
	<-time.After(waitTime)
	c.Destroy(true)
}

func testConn(c Client, t *testing.T) {
	c.Connect(connHandler(t))
}

// conn -> pub
func TestClient_Publish(t *testing.T) {
	c := plainClient()
	testConn(c, t)
	<-time.After(waitTime)
	testPub(c, t)
	<-time.After(waitTime)
	c.Destroy(true)

	c = tlsClient()
	testConn(c, t)
	<-time.After(waitTime)
	testPub(c, t)
	<-time.After(waitTime)
	c.Destroy(true)
}

func testPub(c Client, t *testing.T) {
	c.Publish(func(topic string, code PubAckCode) {
	}, testMsgs...)
}

func TestClient_Subscribe(t *testing.T) {
	c := plainClient()
	testConn(c, t)
	<-time.After(waitTime)
	testSub(c, t)
	<-time.After(waitTime)
	c.Destroy(true)

	c = tlsClient()
	testConn(c, t)
	<-time.After(waitTime)
	testSub(c, t)
	<-time.After(waitTime)
	c.Destroy(true)
}

// conn -> sub -> pub
func testSub(c Client, t *testing.T) {
	count := atomic.Value{}
	count.Store(int(0))
	c.Subscribe(func(topic string, code SubAckCode, msg []byte) {
		if code == SubFail {
			t.FailNow()
		}

		count.Store(count.Load().(int) + 1)
		switch topic {
		case testTopics[0]:
			if code != Qos0 || bytes.Compare(testMsgs[0].Payload, msg) != 0 {
				t.Log("fail at sub topic 0")
				t.FailNow()
			} else {
				t.Log("received sub topic 0")
			}
		case testTopics[1]:
			if code != Qos1 || bytes.Compare(testMsgs[1].Payload, msg) != 0 {
				t.Log("fail at sub topic 1")
				t.FailNow()
			} else {
				t.Log("received sub topic 1")
			}
		case testTopics[2]:
			if code != Qos2 || bytes.Compare(testMsgs[2].Payload, msg) != 0 {
				t.Log("fail at sub topic 2")
				t.FailNow()
			} else {
				t.Log("received sub topic 2")
			}
		}
	}, testSubs...)
	<-time.After(waitTime)
	testPub(c, t)
	<-time.After(waitTime)
	if c := count.Load().(int); c != len(testMsgs) {
		t.Log("sub recv count =", c)
		t.FailNow()
	}
}

// conn -> sub -> pub -> unSub
func TestClient_UnSubscribe(t *testing.T) {
	c := plainClient()
	testUnSub(c, t)
	<-time.After(waitTime)
	c.Destroy(true)

	c = tlsClient()
	testUnSub(c, t)
	<-time.After(waitTime)
	c.Destroy(true)
}

func testUnSub(c Client, t *testing.T) {
	testConn(c, t)
	<-time.After(waitTime)
	testSub(c, t)
	<-time.After(waitTime)
	c.UnSubscribe(func(topic string, code SubAckCode) {}, testTopics...)
	<-time.After(waitTime)
	testPub(c, t)
}
