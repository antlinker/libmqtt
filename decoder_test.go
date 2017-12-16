package libmqtt

import (
	"bufio"
	"bytes"
	"testing"
)

func TestDecodeRemainLength(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.Write([]byte{0x04})
	length, err := decodeRemainLength(buffer)
	if err != nil || length != 0x04 {
		t.Log(length)
		t.Fail()
	}
	buffer.Reset()
}

func TestDecodeOnePacket(t *testing.T) {
	targetBytes := connWillBytes
	buffer := &bytes.Buffer{}
	_, err := buffer.Write(targetBytes)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	reader := bufio.NewReader(buffer)
	pkt, err := decodeOnePacket(reader)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	buffer.Reset()
	switch pkt.(type) {
	case *ConPacket:
		pkt.Bytes(buffer)
		pktBytes := buffer.Bytes()
		if bytes.Compare(pktBytes, targetBytes) != 0 {
			t.Log(pktBytes)
			t.Fail()
		}
	default:
		t.Fail()
	}
}
