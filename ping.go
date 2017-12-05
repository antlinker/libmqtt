package libmqtt

import "bytes"

// PingReqPacket is sent from a Client to the Server.
//
// It can be used to:
// 		1. Indicate to the Server that the Client is alive in the absence of any other Control Packets being sent from the Client to the Server.
// 		2. Request that the Server responds to confirm that it is alive.
// 		3. Exercise the network to indicate that the Network Connection is active.
//
// This Packet is used in Keep Alive processing
type PingReqPacket struct {
}

func (s *PingReqPacket) Type() CtrlType {
	return CtrlUnSubAck
}

func (s *PingReqPacket) Bytes(buffer *bytes.Buffer) (err error) {
	if buffer == nil || s == nil {
		return
	}
	// fixed header
	buffer.WriteByte(CtrlPingReq << 4)
	// remaining length
	return buffer.WriteByte(0x00)
}

// PingRespPacket is sent by the Server to the Client in response to
// a PingReqPacket. It indicates that the Server is alive.
type PingRespPacket struct {
}

func (s *PingRespPacket) Type() CtrlType {
	return CtrlUnSubAck
}

func (s *PingRespPacket) Bytes(buffer *bytes.Buffer) (err error) {
	if buffer == nil || s == nil {
		return
	}
	// fixed header
	buffer.WriteByte(CtrlPingResp << 4)
	return buffer.WriteByte(0x00)
}
