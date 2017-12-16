package libmqtt

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	// test with emqttd server (http://emqtt.io/ or https://github.com/emqtt/emqttd)
	client := NewClient(
		WithServer("localhost:1883"),
		WithDialTimeout(10),
		WithKeepalive(10),
		WithIdentity("admin", "public"),
		WithWill("test", Qos0, false, []byte("test data")),
	)

	client.Connect(func(server string, code ConAckCode) {
		t.Log(code)
		t.Fail()
	})
	client.Wait()
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
		WithKeepalive(10),
		WithCleanSession(),
		WithIdentity("admin", "public"),
		WithWill("test", Qos0, true, []byte("test data")),
	)

	client.Connect(func(server string, code ConAckCode) {
		t.Log(code)
		t.Fail()
	})
	client.Wait()
}

func TestClient_Publish(t *testing.T) {

}

func TestClient_Subscribe(t *testing.T) {

}

func TestClient_UnSubscribe(t *testing.T) {

}
