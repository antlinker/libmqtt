package libmqtt

import (
	"testing"
)

func TestPingReqPacket_Bytes(t *testing.T) {
	testBytes(testPingReqMsg, testPingReqMsgBytes, t)
}

func TestPingRespPacket_Bytes(t *testing.T) {
	testBytes(testPingRespMsg, testPingRespMsgBytes, t)
}
