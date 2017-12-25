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
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	// test with emqttd server (http://emqtt.io/ or https://github.com/emqtt/emqttd)
	client := NewClient(
		WithServer("localhost:1883"),
		WithDialTimeout(10),
		WithKeepalive(10, 1.2),
		WithIdentity("admin", "public"),
		WithWill("test", Qos0, false, []byte("test data")),
		WithLogger(Verbose),
	)

	client.Connect(func(server string, code ConAckCode) {
		t.Log(code)
		t.Fail()
	})

	<-time.After(5 * time.Second)
	client.Destroy(true)
}

func TestNewClientTLS(t *testing.T) {
	// test with emqttd server (http://emqtt.io/ or https://github.com/emqtt/emqttd)
	// using server default cert and key pair
	client := NewClient(
		WithServer("localhost:8883"),
		WithTLS(
			"./testdata/client-cert.pem",
			"./testdata/client-key.pem",
			"./testdata/ca-cert.pem",
			"MacBook-Air.local",
			true),
		WithDialTimeout(10),
		WithKeepalive(10, 1.2),
		WithCleanSession(),
		WithIdentity("admin", "public"),
		WithWill("test", Qos0, true, []byte("test data")),
		WithLogger(Verbose),
	)

	client.Connect(func(server string, code ConAckCode) {
		t.Log(code)
		t.Fail()
	})

	<-time.After(5 * time.Second)
	client.Destroy(true)
}

func TestClient_Publish(t *testing.T) {

}

func TestClient_Subscribe(t *testing.T) {

}

func TestClient_UnSubscribe(t *testing.T) {

}
