package libmqtt

import (
	"bufio"
	"bytes"
	"math"
	"testing"

	std "github.com/eclipse/paho.mqtt.golang/packets"
)

// common test data
var (
	testUsername     = "foo"
	testPassword     = "bar"
	testClientID     = "1"
	testCleanSession = true
	testWill         = true
	testWillQos      = Qos2
	testWillRetain   = true
	testWillTopic    = "foo"
	testWillMessage  = []byte("bar")
	testKeepalive    = uint16(1)
	testProtoVersion = V311
	testPubDup       = false

	testTopics = []string{
		"/test",
		"/test/foo",
		"/test/bar",
	}
	testTopicQos = []QosLevel{
		Qos0,
		Qos1,
		Qos2,
	}
	testTopicMsgs = []string{
		"test data qos 0",
		"foo data qos 1",
		"bar data qos 2",
	}
	testPacketID uint16 = math.MaxUint16 / 2

	testConnackPresent = true
	testConnackCode    = ConnAccepted

	testSubAckCodes = []SubAckCode{SubOkMaxQos0, SubOkMaxQos1, SubFail}
)

// conn test data
var (
	testConnWillMsg      *ConnPacket
	testConnWillMsgBytes []byte

	testConnMsg      *ConnPacket
	testConnMsgBytes []byte

	testConnAckMsg      *ConnAckPacket
	testConnAckMsgBytes []byte

	testDisConnMsg      = DisConnPacket
	testDisConnMsgBytes []byte
)

func initConn() {
	testConnWillMsg = &ConnPacket{
		Username:     testUsername,
		Password:     testPassword,
		protoLevel:   testProtoVersion,
		ClientID:     testClientID,
		CleanSession: testCleanSession,
		IsWill:       testWill,
		WillQos:      testWillQos,
		WillRetain:   testWillRetain,
		WillTopic:    testWillTopic,
		WillMessage:  testWillMessage,
		Keepalive:    testKeepalive,
	}
	connWillPkt := std.NewControlPacket(std.Connect).(*std.ConnectPacket)
	connWillPkt.Username = testUsername
	connWillPkt.UsernameFlag = true
	connWillPkt.Password = []byte(testPassword)
	connWillPkt.PasswordFlag = true
	connWillPkt.ProtocolName = "MQTT"
	connWillPkt.ProtocolVersion = testProtoVersion
	connWillPkt.ClientIdentifier = testClientID
	connWillPkt.CleanSession = testCleanSession
	connWillPkt.WillFlag = testWill
	connWillPkt.WillQos = testWillQos
	connWillPkt.WillRetain = testWillRetain
	connWillPkt.WillTopic = testWillTopic
	connWillPkt.WillMessage = testWillMessage
	connWillPkt.Keepalive = testKeepalive
	connWillBuf := &bytes.Buffer{}
	connWillPkt.Write(connWillBuf)
	testConnWillMsgBytes = connWillBuf.Bytes()

	testConnMsg = &ConnPacket{
		Username:     testUsername,
		Password:     testPassword,
		protoLevel:   testProtoVersion,
		ClientID:     testClientID,
		CleanSession: testCleanSession,
		Keepalive:    testKeepalive,
	}
	connPkt := std.NewControlPacket(std.Connect).(*std.ConnectPacket)
	connPkt.Username = testUsername
	connPkt.UsernameFlag = true
	connPkt.Password = []byte(testPassword)
	connPkt.PasswordFlag = true
	connPkt.ProtocolName = "MQTT"
	connPkt.ProtocolVersion = testProtoVersion
	connPkt.ClientIdentifier = testClientID
	connPkt.CleanSession = testCleanSession
	connPkt.Keepalive = testKeepalive
	connBuf := &bytes.Buffer{}
	connPkt.Write(connBuf)
	testConnMsgBytes = connBuf.Bytes()

	testConnAckMsg = &ConnAckPacket{
		Present: testConnackPresent,
		Code:    testConnackCode,
	}
	connackPkt := std.NewControlPacket(std.Connack).(*std.ConnackPacket)
	connackPkt.SessionPresent = testConnackPresent
	connackPkt.ReturnCode = testConnackCode
	connAckBuf := &bytes.Buffer{}
	connackPkt.Write(connAckBuf)
	testConnAckMsgBytes = connAckBuf.Bytes()

	disConnPkt := std.NewControlPacket(std.Disconnect).(*std.DisconnectPacket)
	disconnBuf := &bytes.Buffer{}
	disConnPkt.Write(disconnBuf)
	testDisConnMsgBytes = disconnBuf.Bytes()
}

// sub test data
var (
	testSubTopics []*Topic

	testSubMsgs     []*SubscribePacket
	testSubMsgBytes [][]byte

	testSubAckMsgs     []*SubAckPacket
	testSubAckMsgBytes [][]byte

	testUnSubMsgs     []*UnSubPacket
	testUnSubMsgBytes [][]byte

	testUnSubAckMsg      *UnSubAckPacket
	testUnSubAckMsgBytes []byte
)

func initSub() {
	testSubTopics = make([]*Topic, len(testTopics))
	for i := range testSubTopics {
		testSubTopics[i] = &Topic{Name: testTopics[i], Qos: testTopicQos[i]}
	}

	testSubMsgs = make([]*SubscribePacket, len(testTopics))
	testSubMsgBytes = make([][]byte, len(testTopics))
	testSubAckMsgs = make([]*SubAckPacket, len(testTopics))
	testSubAckMsgBytes = make([][]byte, len(testTopics))
	testUnSubMsgs = make([]*UnSubPacket, len(testTopics))
	testUnSubMsgBytes = make([][]byte, len(testTopics))

	for i := range testTopics {
		testSubMsgs[i] = &SubscribePacket{
			Topics:   testSubTopics[:i+1],
			PacketID: testPacketID,
		}

		subPkt := std.NewControlPacket(std.Subscribe).(*std.SubscribePacket)
		subPkt.Topics = testTopics[:i+1]
		subPkt.Qoss = testTopicQos[:i+1]
		subPkt.MessageID = testPacketID

		subBuf := &bytes.Buffer{}
		subPkt.Write(subBuf)
		testSubMsgBytes[i] = subBuf.Bytes()

		testSubAckMsgs[i] = &SubAckPacket{
			PacketID: testPacketID,
			Codes:    testSubAckCodes[:i+1],
		}
		subAckPkt := std.NewControlPacket(std.Suback).(*std.SubackPacket)
		subAckPkt.MessageID = testPacketID
		subAckPkt.ReturnCodes = testSubAckCodes[:i+1]
		subAckBuf := &bytes.Buffer{}
		subAckPkt.Write(subAckBuf)
		testSubAckMsgBytes[i] = subAckBuf.Bytes()

		testUnSubMsgs[i] = &UnSubPacket{
			PacketID:   testPacketID,
			TopicNames: testTopics[:i+1],
		}
		unsubPkt := std.NewControlPacket(std.Unsubscribe).(*std.UnsubscribePacket)
		unsubPkt.Topics = testTopics[:i+1]
		unsubPkt.MessageID = testPacketID
		unSubBuf := &bytes.Buffer{}
		unsubPkt.Write(unSubBuf)
		testUnSubMsgBytes[i] = unSubBuf.Bytes()
	}

	unSunAckBuf := &bytes.Buffer{}
	testUnSubAckMsg = &UnSubAckPacket{PacketID: testPacketID}
	unsubAckPkt := std.NewControlPacket(std.Unsuback).(*std.UnsubackPacket)
	unsubAckPkt.MessageID = testPacketID
	unsubAckPkt.Write(unSunAckBuf)
	testUnSubAckMsgBytes = unSunAckBuf.Bytes()
}

// pub test data
var (
	testPubMsgs     []*PublishPacket
	testPubMsgBytes [][]byte

	testPubAckMsg      *PubAckPacket
	testPubAckMsgBytes []byte

	testPubRecvMsg      *PubRecvPacket
	testPubRecvMsgBytes []byte

	testPubRelMsg      *PubRelPacket
	testPubRelMsgBytes []byte

	testPubCompMsg      *PubCompPacket
	testPubCompMsgBytes []byte
)

// init pub test data
func initPub() {
	testPubMsgs = make([]*PublishPacket, len(testTopics))
	testPubMsgBytes = make([][]byte, len(testTopics))
	for i := range testTopics {
		testPubMsgs[i] = &PublishPacket{
			IsDup:     testPubDup,
			TopicName: testTopics[i],
			Qos:       testTopicQos[i],
			Payload:   []byte(testTopicMsgs[i]),
			PacketID:  testPacketID,
		}

		// create standard publish packet and make bytes
		pkt := std.NewControlPacketWithHeader(std.FixedHeader{
			MessageType: std.Publish,
			Dup:         testPubDup,
			Qos:         testTopicQos[i],
		}).(*std.PublishPacket)
		pkt.TopicName = testTopics[i]
		pkt.Payload = []byte(testTopicMsgs[i])
		pkt.MessageID = testPacketID

		buf := &bytes.Buffer{}
		pkt.Write(buf)
		testPubMsgBytes[i] = buf.Bytes()
	}

	testPubAckMsg = &PubAckPacket{PacketID: testPacketID}
	pubAckPkt := std.NewControlPacket(std.Puback).(*std.PubackPacket)
	pubAckPkt.MessageID = testPacketID
	pubAckBuf := &bytes.Buffer{}
	pubAckPkt.Write(pubAckBuf)
	testPubAckMsgBytes = pubAckBuf.Bytes()

	pubRecvBuf := &bytes.Buffer{}
	testPubRecvMsg = &PubRecvPacket{PacketID: testPacketID}
	pubRecPkt := std.NewControlPacket(std.Pubrec).(*std.PubrecPacket)
	pubRecPkt.MessageID = testPacketID
	pubRecPkt.Write(pubRecvBuf)
	testPubRecvMsgBytes = pubRecvBuf.Bytes()

	pubRelBuf := &bytes.Buffer{}
	testPubRelMsg = &PubRelPacket{PacketID: testPacketID}
	pubRelPkt := std.NewControlPacket(std.Pubrel).(*std.PubrelPacket)
	pubRelPkt.MessageID = testPacketID
	pubRelPkt.Write(pubRelBuf)
	testPubRelMsgBytes = pubRelBuf.Bytes()

	pubCompBuf := &bytes.Buffer{}
	testPubCompMsg = &PubCompPacket{PacketID: testPacketID}
	pubCompPkt := std.NewControlPacket(std.Pubcomp).(*std.PubcompPacket)
	pubCompPkt.MessageID = testPacketID
	pubCompPkt.Write(pubCompBuf)
	testPubCompMsgBytes = pubCompBuf.Bytes()
}

var (
	testPingReqMsg      = PingReqPacket
	testPingReqMsgBytes []byte

	testPingRespMsg      = PingRespPacket
	testPingRespMsgBytes []byte
)

func initPing() {
	reqBuf := &bytes.Buffer{}
	req := std.NewControlPacket(std.Pingreq).(*std.PingreqPacket)
	req.Write(reqBuf)
	testPingReqMsgBytes = reqBuf.Bytes()

	respBuf := &bytes.Buffer{}
	resp := std.NewControlPacket(std.Pingresp).(*std.PingrespPacket)
	resp.Write(respBuf)
	testPingRespMsgBytes = respBuf.Bytes()
}

func testBytes(pkt Packet, target []byte, t *testing.T) {
	buf := &bytes.Buffer{}
	w := bufio.NewWriter(buf)
	if err := pkt.Bytes(w); err != nil {
		t.Log(err)
		t.Fail()
	}
	if bytes.Compare(buf.Bytes(), target) != 0 {
		t.Log("S", buf.Bytes())
		t.Log("T", target)
		t.Fail()
	}
}

func init() {
	initPing()
	initConn()
	initSub()
	initPub()
}
