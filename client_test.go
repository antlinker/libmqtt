package libmqtt

import "testing"

func TestNewClientTLS(t *testing.T) {
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

	client.Connect(func(server string, code ConnAckCode) {
		t.Log(code)
		t.Fail()
	})
	client.Wait()
}

func TestNewClient(t *testing.T) {
	client := NewClient(
		WithServer("localhost:1883"),
		WithDialTimeout(10),
		WithKeepalive(10),
		WithIdentity("admin", "public"),
		WithWill("test", Qos0, false, []byte("test data")),
	)

	client.Connect(func(server string, code ConnAckCode) {
		t.Log(code)
		t.Fail()
	})
	client.Wait()
}
