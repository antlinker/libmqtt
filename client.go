/*
 * Copyright GoIIoT (https://github.com/goiiot)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package libmqtt

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"math"
	"net"
	"sync"
	"time"
)

var lg *logger

// BackoffOption defines the parameters for the reconnecting backoff strategy.
type BackoffOption struct {
	// MaxDelay defines the upper bound of backoff delay,
	// which is time in second.
	MaxDelay uint16

	// FirstDelay is the time to wait before retrying after the first failure,
	// also time in second.
	FirstDelay uint16

	// Factor is applied to the backoff after each retry.
	// e.g. FirstDelay = 1 and Factor = 2, then the SecondDelay is 2, the ThirdDelay is 4s
	Factor float32
}

// Option is client option for connection options
type Option func(*client)

// WithCleanSession will set clean flag in connect packet
func WithCleanSession() Option {
	return func(c *client) {
		c.options.cleanSession = true
	}
}

// WithIdentity for username and password
func WithIdentity(username, password string) Option {
	return func(c *client) {
		c.options.username = username
		c.options.password = password
	}
}

// WithKeepalive set the keepalive interval (time in second)
func WithKeepalive(keepalive uint16, factor float64) Option {
	return func(c *client) {
		c.options.keepalive = time.Duration(keepalive) * time.Second
		if factor > 1 {
			c.options.keepaliveFactor = factor
		} else {
			factor = 1.2
		}
	}
}

// WithBackoffStrategy will set reconnect backoff strategy
func WithBackoffStrategy(bf *BackoffOption) Option {
	return func(c *client) {
		if bf != nil {
			c.options.bf = bf
		}
	}
}

// WithClientID set the client id for connection
func WithClientID(clientID string) Option {
	return func(c *client) {
		c.options.clientID = clientID
	}
}

// WithWill mark this connection as a will teller
func WithWill(topic string, qos QosLevel, retain bool, payload []byte) Option {
	return func(c *client) {
		c.options.isWill = true
		c.options.willTopic = topic
		c.options.willQos = qos
		c.options.willRetain = retain
		c.options.willPayload = payload
	}
}

// WithServer adds servers as client server
// Just use "ip:port" or "domain.name:port"
// Only TCP connection supported for now
func WithServer(servers ...string) Option {
	return func(c *client) {
		c.options.servers = servers
	}
}

// WithTLS for client tls certification
func WithTLS(certFile, keyFile string, caCert string, serverNameOverride string, skipVerify bool) Option {
	return func(c *client) {
		b, err := ioutil.ReadFile(caCert)
		if err != nil {
			panic("load ca cert file failed ")
		}
		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM(b) {
			panic("append certificates failed ")
		}
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			panic("load client cert file failed ")
		}

		c.options.tlsConfig = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: skipVerify,
			ClientCAs:          cp,
			ServerName:         serverNameOverride,
		}
	}
}

// WithDialTimeout for connection time out (time in second)
func WithDialTimeout(timeout uint16) Option {
	return func(c *client) {
		c.options.dialTimeout = time.Duration(timeout) * time.Second
	}
}

// WithSendBuf designate the channel size of send
func WithSendBuf(size int) Option {
	return func(c *client) {
		if size < 1 {
			size = 1
		} else if size > 1024 {
			size = 1024
		}
		c.options.sendChanSize = size
	}
}

// WithRecvBuf designate the channel size of receive
func WithRecvBuf(size int) Option {
	return func(c *client) {
		if size < 1 {
			size = 1
		} else if size > 1024 {
			size = 1024
		}
		c.options.recvChanSize = size
	}
}

// WithRouter set the router for topic dispatch
func WithRouter(r TopicRouter) Option {
	return func(c *client) {
		if r != nil {
			c.router = r
		}
	}
}

// WithLogger will create basic logger for log
func WithLogger(l LogLevel) Option {
	return func(c *client) {
		lg = newLogger(l)
	}
}

// NewClient will create a new mqtt client
func NewClient(options ...Option) Client {
	c := defaultClient()
	for _, o := range options {
		o(c)
	}

	c.sendC = make(chan Packet, c.options.sendChanSize)
	c.recvC = make(chan PublishPacket, c.options.recvChanSize)
	return c
}

// clientOptions is the options for client to connect, reconnect, disconnect
type clientOptions struct {
	sendChanSize    int            // send channel size
	recvChanSize    int            // recv channel size
	servers         []string       // server address strings
	dialTimeout     time.Duration  // dial timeout in second
	clientID        string         // used by ConPacket
	username        string         // used by ConPacket
	password        string         // used by ConPacket
	keepalive       time.Duration  // used by ConPacket (time in second)
	keepaliveFactor float64        // used for reasonable amount time to close conn if no ping resp
	cleanSession    bool           // used by ConPacket
	isWill          bool           // used by ConPacket
	willTopic       string         // used by ConPacket
	willPayload     []byte         // used by ConPacket
	willQos         byte           // used by ConPacket
	willRetain      bool           // used by ConPacket
	tlsConfig       *tls.Config    // tls config with client side cert
	bf              *BackoffOption // backoff option for client reconnection
}

// Client act as a mqtt client
type Client interface {
	// Connect to all specified server with client options
	Connect(h ConnHandler)

	// Publish a message for the topic
	Publish(msg ...*PublishPacket)

	// Handle register topic handlers, mostly used for RegexHandler, RestHandler
	Handle(topic string, h SubHandler)

	// Subscribe topic(s)
	Subscribe(topics ...*Topic)

	// UnSubscribe topic(s)
	UnSubscribe(topics ...string)

	// Wait will wait until all connection finished
	Wait()

	// Destroy all client connection
	Destroy(force bool)
}

type client struct {
	options *clientOptions     // client connection options
	subs    *sync.Map          // Topic(s) -> []SubHandler
	conn    *sync.Map          // ServerAddr -> connection
	sendC   chan Packet        // Pub channel for sending publish packet to server
	recvC   chan PublishPacket // Pub recv channel for receiving
	idGen   *idGenerator       // sorted in use packetId []uint16
	router  TopicRouter        // topic router
	workers *sync.WaitGroup    // workers
}

// defaultClient create the client with default options
func defaultClient() *client {
	c := &client{
		options: &clientOptions{
			sendChanSize: 128,
			recvChanSize: 128,
			bf: &BackoffOption{
				MaxDelay:   120, // default max retry delay is 2min
				FirstDelay: 1,   // first retry delay is 1s
				Factor:     1.5,
			},
			dialTimeout:     20 * time.Second, // default timeout when dial to server
			keepalive:       2 * time.Minute,  // default keepalive interval is 2min
			keepaliveFactor: 1.5,              // default reasonable amount of time 3min
		},
		router:  &TextRouter{},
		subs:    &sync.Map{},
		conn:    &sync.Map{},
		idGen:   newIDGenerator(),
		workers: &sync.WaitGroup{},
	}
	return c
}

// Connect to all designated server
func (c *client) Connect(h ConnHandler) {
	for _, s := range c.options.servers {
		c.workers.Add(1)
		go c.connect(s, h)
	}

	go func() {
		for pkt := range c.recvC {
			c.router.Dispatch(pkt)
		}
	}()
}

// Publish message(s) to topic(s), one to one
func (c *client) Publish(msg ...*PublishPacket) {
	for _, m := range msg {
		if m.Qos > Qos2 {
			panic("Invalid QoS level, should either be QoS0, QoS1 or QoS2 ")
		}

		toSend := &PublishPacket{
			Qos:       m.Qos,
			IsRetain:  m.IsRetain,
			TopicName: m.TopicName,
			Payload:   m.Payload,
		}

		if toSend.Qos != Qos0 {
			toSend.PacketID = c.nextID()
		}
		c.sendC <- toSend
	}
}

// Handle subscription message route
func (c *client) Handle(topic string, h SubHandler) {
	if h != nil {
		lg.d("HANDLE registered handler, topic =", topic)
		c.router.Handle(topic, h)
	}
}

// SubScribe topic(s)
func (c *client) Subscribe(topics ...*Topic) {
	// send sub message
	lg.d("SEND subscribe, topic(s) =", topics)
	c.sendC <- &SubscribePacket{
		Topics:   topics,
		PacketID: c.nextID(),
	}
}

// UnSubscribe topic(s)
func (c *client) UnSubscribe(topics ...string) {
	for _, t := range topics {
		c.subs.Delete(t)
	}

	lg.d("SEND UnSub, topic(s) =", topics)
	c.sendC <- &UnSubPacket{
		TopicNames: topics,
		PacketID:   c.nextID(),
	}
}

// Wait will wait for all connection to exit
func (c *client) Wait() {
	c.workers.Wait()
}

// Destroy will disconnect form all server
// If force is true, then close connection without sending a DisConnPacket
func (c *client) Destroy(force bool) {
	if force {
		c.conn.Range(func(k, v interface{}) bool {
			va := v.(*connImpl)
			va.close()
			return true
		})
	} else {

	}
}

// connect to one server and start mqtt logic
func (c *client) connect(server string, h ConnHandler) {
	defer c.workers.Done()

	var conn net.Conn
	var err error

	if c.options.tlsConfig != nil {
		// with tls
		conn, err = tls.DialWithDialer(&net.Dialer{Timeout: c.options.dialTimeout}, "tcp", server, c.options.tlsConfig)
		if err != nil {
			lg.e("connection with tls failed", err)
			h(server, math.MaxUint8, err)
			return
		}
	} else {
		// without tls
		conn, err = net.DialTimeout("tcp", server, c.options.dialTimeout)
		if err != nil {
			lg.e("connection failed", err)
			h(server, math.MaxUint8, err)
			return
		}
	}

	connImpl := &connImpl{
		parent:     c,
		name:       server,
		conn:       conn,
		sendBuf:    &bytes.Buffer{},
		keepaliveC: make(chan int),
		recvC:      make(chan Packet),
	}

	go connImpl.handleSend()
	go connImpl.handleRecv()

	connImpl.send(&ConPacket{
		Username:     c.options.username,
		Password:     c.options.password,
		ClientID:     c.options.clientID,
		CleanSession: c.options.cleanSession,
		IsWill:       c.options.isWill,
		WillQos:      c.options.willQos,
		WillTopic:    c.options.willTopic,
		WillMessage:  c.options.willPayload,
		WillRetain:   c.options.willRetain,
		Keepalive:    uint16(c.options.keepalive / time.Second),
	})

	select {
	case pkt, more := <-connImpl.recvC:
		if more {
			if pkt.Type() == CtrlConnAck {
				p := pkt.(*ConAckPacket)
				if p.Code != ConnAccepted {
					h(server, p.Code, nil)
					return
				}
			} else {
				h(server, math.MaxUint8, ErrBadPacket)
				return
			}
		} else {
			h(server, math.MaxUint8, ErrBadPacket)
			return
		}
	case <-time.After(c.options.dialTimeout):
		h(server, math.MaxUint8, ErrTimeOut)
		return
	}

	lg.i("CONN connection success, server =", server)
	go h(server, ConnAccepted, nil)

	// login success, start mqtt logic
	c.conn.Store(server, connImpl)
	connImpl.start()
}

// get next valid packet id
func (c *client) nextID() uint16 {
	return c.idGen.next()
}

// free packet id
func (c *client) freeID(id uint16) {
	c.idGen.free(id)
}

// connImpl is the wrapper of connection to server
// tend to actual packet send and receive
type connImpl struct {
	parent     *client       // client which created this connection
	name       string        // server addr info
	conn       net.Conn      // connection to server
	sendBuf    *bytes.Buffer // buffer for packet send
	recvC      chan Packet   // received packet from server
	keepaliveC chan int      // keepalive packet
}

// start mqtt logic
func (c *connImpl) start() {
	// start keepalive if required
	if c.parent.options.keepalive > 0 {
		go c.keepalive()
	}

	// inspect incoming packet
	for pkt := range c.recvC {
		switch pkt.Type() {
		case CtrlSubAck:
			p := pkt.(*SubAckPacket)
			lg.d("RECV SubAck, id =", p.PacketID)
			c.parent.freeID(p.PacketID)
			// TODO: notify Sub QoS response
		case CtrlUnSubAck:
			p := pkt.(*UnSubAckPacket)
			lg.d("RECV UnSubAck, id =", p.PacketID)
			c.parent.freeID(p.PacketID)
		case CtrlPublish:
			p := pkt.(*PublishPacket)
			lg.d("RECV Publish, id =", p.PacketID, "QoS =", p.Qos)
			c.parent.recvC <- *p

			// tend to QoS
			switch p.Qos {
			case Qos2:
				lg.d("SEND PubAck for Publish, id =", p.PacketID)
				c.send(&PubAckPacket{PacketID: p.PacketID})
			case Qos1:
				lg.d("SEND PubAck for Publish, id =", p.PacketID)
				c.send(&PubRecvPacket{PacketID: p.PacketID})
			}
		case CtrlPubAck:
			p := pkt.(*PubAckPacket)
			lg.d("RECV PubAck, id =", p.PacketID)

			c.parent.freeID(p.PacketID)
		case CtrlPubRecv:
			p := pkt.(*PubRecvPacket)
			lg.d("RECV PubRec, id =", p.PacketID)

			c.send(&PubRelPacket{PacketID: p.PacketID})
			lg.d("SEND PubRel, id =", p.PacketID)
		case CtrlPubRel:
			p := pkt.(*PubRelPacket)
			lg.d("RECV PubRel, id =", p.PacketID)

			c.send(&PubCompPacket{PacketID: p.PacketID})
			lg.d("SEND PubComp, id =", p.PacketID)
		case CtrlPubComp:
			p := pkt.(*PubCompPacket)
			lg.d("RECV PubComp id =", p.PacketID)

			c.parent.freeID(p.PacketID)
		default:
			lg.d("RECV packet, type =", pkt.Type())
		}
	}
}

// keepalive with server
func (c *connImpl) keepalive() {
	lg.d("START keepalive")
	defer lg.d("END keepalive")

	t := time.NewTicker(c.parent.options.keepalive)
	defer t.Stop()

	for range t.C {
		c.send(PingReqPacket)

		select {
		case _, more := <-c.keepaliveC:
			if !more {
				return
			}
		case <-time.After(c.parent.options.keepalive * time.Duration(c.parent.options.keepaliveFactor)):
			// ping timeout
			t.Stop()
			c.conn.Close()
			return
		}
	}
}

// close this connection
func (c *connImpl) close() {
	lg.i("END connection, server =", c.name)
	c.send(DisConPacket)
	c.conn.Close()
}

// handle client message send
func (c *connImpl) handleSend() {
	for pkt := range c.parent.sendC {
		pkt.Bytes(c.sendBuf)
		if _, err := c.sendBuf.WriteTo(c.conn); err != nil {
			break
		}

		if pkt.Type() == CtrlDisConn {
			// disconnect to server
			break
		}
	}
}

// handle all message recv
func (c *connImpl) handleRecv() {
	for {
		pkt, err := decodeOnePacket(c.conn)
		if err != nil {
			lg.e("CONN broken", "server =", c.name, "err =", err)
			close(c.recvC)
			close(c.keepaliveC)
			break
		}

		// pass packets
		if pkt == PingRespPacket {
			lg.d("RECV keepalive message")
			c.keepaliveC <- 1
		} else {
			c.recvC <- pkt
		}
	}
}

// send internal mqtt logic packet
func (c *connImpl) send(pkt Packet) {
	c.parent.sendC <- pkt
}
