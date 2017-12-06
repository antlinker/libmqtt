package libmqtt

import (
	"bytes"
	"testing"
)

var (
	connWillPkg = &ConnPacket{
		Username:   "user",
		Password:   "pass",
		protoLevel: V311,

		ClientId:     "1",
		CleanSession: true,
		IsWill:       true,
		WillQos:      Qos2,
		WillRetain:   true,
		WillTopic:    "lost",
		WillMessage:  []byte("peace"),
		Keepalive:    10,
	}
	connPkg = &ConnPacket{
		Username:   "user",
		Password:   "pass",
		protoLevel: V311,

		ClientId:     "1",
		CleanSession: true,
		Keepalive:    10,
	}
	connWillBytes = []byte{
		0x10,                 // fixed header: conn:0
		38,                   // remaining length: 38
		0, 4, 77, 81, 84, 84, // Protocol Name: "MQTT"
		4,     // Protocol Level 3.1.1
		0xF6,  // connect flags: 11110110
		0, 10, // keepalive: 10s
		0, 1, 49, // client id: "1"
		0, 4, 108, 111, 115, 116, // will topic: "lost"
		0, 5, 112, 101, 97, 99, 101, // will msg: "peace"
		0, 4, 117, 115, 101, 114, // Username: "user"
		0, 4, 112, 97, 115, 115, // Password: "pass"
	}
	connBytes = []byte{
		0x10,                 // fixed header: conn:0
		25,                   // remaining length: 38
		0, 4, 77, 81, 84, 84, // Protocol Name: "MQTT"
		4,     // Protocol Level 3.1.1
		0xC2,  // connect flags: 11000010
		0, 10, // keepalive: 10s
		0, 1, 49, // client id: "1"
		0, 4, 117, 115, 101, 114, // Username: "user"
		0, 4, 112, 97, 115, 115, // Password: "pass"
	}
)

func BenchmarkConnPacket_Bytes(b *testing.B) {
	buffer := &bytes.Buffer{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		connWillPkg.Bytes(buffer)
	}
}

func TestConnPacket(t *testing.T) {
	buffer := &bytes.Buffer{}
	connWillPkg.Bytes(buffer)
	testConnWillBytes := buffer.Bytes()
	if bytes.Compare(testConnWillBytes, connWillBytes) != 0 {
		t.Log(testConnWillBytes)
		t.Fail()
	}

	buffer.Reset()
	connPkg.Bytes(buffer)
	testConnBytes := buffer.Bytes()
	if bytes.Compare(testConnBytes, connBytes) != 0 {
		t.Log(testConnBytes)
		t.Fail()
	}
}

var (
	connAckPkg   = &ConnAckPacket{Present: true, Code: ConnAccepted}
	connAckBytes = []byte{0x20, 0x02, 0x01, 0x00}
)

func BenchmarkConnAckPacket_Bytes(b *testing.B) {
	buffer := &bytes.Buffer{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		connAckPkg.Bytes(buffer)
	}
}

func TestDisConnAckPacket(t *testing.T) {
	buffer := &bytes.Buffer{}
	connAckPkg.Bytes(buffer)
	testConnAckBytes := buffer.Bytes()
	if bytes.Compare(testConnAckBytes, connAckBytes) != 0 {
		t.Fail()
	}
}

var (
	disConnPkg = &DisConnPacket{}
)

func BenchmarkDisConnPacket_Bytes(b *testing.B) {
	buffer := &bytes.Buffer{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		disConnPkg.Bytes(buffer)
	}
}

func TestDisConnPacket(t *testing.T) {
	buffer := &bytes.Buffer{}
	disConnPkg.Bytes(buffer)
	testDisConnBytes := buffer.Bytes()
	if bytes.Compare(testDisConnBytes, disConnBytes) != 0 {
		t.Fail()
	}
}
