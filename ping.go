package libmqtt

import "bytes"

var (
	PingReqPacket  = &pingReqPacket{}
	PingRespPacket = &pingRespPacket{}
)

// pingReqPacket is sent from a Client to the Server.
//
// It can be used to:
// 		1. Indicate to the Server that the Client is alive in the absence of any other Control Packets being sent from the Client to the Server.
// 		2. Request that the Server responds to confirm that it is alive.
// 		3. Exercise the network to indicate that the Network Connection is active.
//
// This Packet is used in Keep Alive processing
type pingReqPacket struct {
}

func (s *pingReqPacket) Type() CtrlType {
	return CtrlUnSubAck
}

func (s *pingReqPacket) Bytes(buffer *bytes.Buffer) (err error) {
	if buffer == nil || s == nil {
		return
	}
	// fixed header
	buffer.WriteByte(CtrlPingReq << 4)
	// remaining length
	return buffer.WriteByte(0x00)
}

// pingRespPacket is sent by the Server to the Client in response to
// a pingReqPacket. It indicates that the Server is alive.
type pingRespPacket struct {
}

func (s *pingRespPacket) Type() CtrlType {
	return CtrlUnSubAck
}

func (s *pingRespPacket) Bytes(buffer *bytes.Buffer) (err error) {
	if buffer == nil || s == nil {
		return
	}

	_, err = buffer.Write(pingRespBytes)
	return
}
