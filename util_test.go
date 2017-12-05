package libmqtt

import "testing"

func BenchmarkBoolToByte(b *testing.B) {
	for i := 0; i < b.N; i++ {
		boolToByte(true)
	}
}

func TestBoolToByte(t *testing.T) {
	if boolToByte(false) != 0x00 {
		t.Fail()
	}

	if boolToByte(true) != 0x01 {
		t.Fail()
	}
}
