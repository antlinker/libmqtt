package libmqtt

const (
	maxMsgSize = 0xffffff7f
)

type CtrlType = byte

const (
	CtrlConn      CtrlType = iota + 1
	CtrlConnAck
	CtrlPublish
	CtrlPubAck
	CtrlPubRecv
	CtrlPubRel
	CtrlPubComp
	CtrlSubscribe
	CtrlSubAck
	CtrlUnSub
	CtrlUnSubAck
	CtrlPingReq
	CtrlPingResp
	CtrlDisConn
)

type CtrlFlag uint8
