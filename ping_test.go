package libmqtt

import (
	"bytes"
	"testing"
)

var (
	pingReqPkg   = &PingReqPacket{}
	pingReqBytes = []byte{0xC0, 0x00}
)

func BenchmarkPingReqPacket_Bytes(b *testing.B) {
	buffer := &bytes.Buffer{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pingReqPkg.Bytes(buffer)
	}
}

func TestPingReqPacket_Bytes(t *testing.T) {
	buffer := &bytes.Buffer{}
	pingReqPkg.Bytes(buffer)
	testPingReqBytes := buffer.Bytes()
	if bytes.Compare(testPingReqBytes, pingReqBytes) != 0 {
		t.Fail()
	}
}

var (
	pingRespPkg   = &PingRespPacket{}
	pingRespBytes = []byte{0xD0, 0x00}
)

func BenchmarkPingRespPacket_Bytes(b *testing.B) {
	buffer := &bytes.Buffer{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pingRespPkg.Bytes(buffer)
	}
}

func TestPingRespPacket_Bytes(t *testing.T) {
	buffer := &bytes.Buffer{}
	pingRespPkg.Bytes(buffer)
	testPingRespBytes := buffer.Bytes()
	if bytes.Compare(testPingRespBytes, pingRespBytes) != 0 {
		t.Fail()
	}
}
