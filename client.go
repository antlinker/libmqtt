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
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var (
	lg *logger
)

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

// WithClientId set the client id for connection
func WithClientId(clientId string) Option {
	return func(c *client) {
		c.options.clientId = clientId
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
	return c
}

// clientOptions is the options for client to connect, reconnect, disconnect
type clientOptions struct {
	sendChanSize    int            // send channel size
	recvChanSize    int            // recv channel size
	servers         []string       // server address strings
	dialTimeout     time.Duration  // dial timeout in second
	clientId        string         // used by ConPacket
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
	Connect(h ConHandler)

	// Publish a message for the topic
	Publish(h PubHandler, msg ...*TopicMsg)

	// Subscribe topic(s)
	Subscribe(h SubHandler, topics ...*Topic)

	// UnSubscribe topic(s)
	UnSubscribe(h UnSubHandler, topics ...string)

	// Wait will wait until all connection finished
	Wait()

	// Destroy all client connection
	Destroy(force bool)
}

type clientAction func()

type client struct {
	options  *clientOptions
	actC     chan clientAction       // action channel
	subC     chan PublishPacket      // sub channel
	pubC     chan PublishPacket      // pub channel
	subs     map[string][]SubHandler // topic -> []SubHandler
	conn     *sync.Map               // server -> connection
	packetId *atomic.Value           // current packet id
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
		actC:     make(chan clientAction),
		pubC:     make(chan PublishPacket, 256),
		subC:     make(chan PublishPacket, 256),
		subs:     make(map[string][]SubHandler),
		conn:     &sync.Map{},
		packetId: &atomic.Value{},
	}
	c.packetId.Store(uint16(0))
	return c
}

// Connect to all designated server
func (c *client) Connect(h ConHandler) {
	for _, s := range c.options.servers {
		go c.connect(s, h)
	}
}

// Publish message(s) to topic(s), one to one
func (c *client) Publish(h PubHandler, msg ...*TopicMsg) {
	for _, m := range msg {
		if m.Qos > Qos2 {
			panic("invalid QoS level, should either be 0, 1 or 2 ")
		}

		c.pubC <- PublishPacket{
			Qos:       m.Qos,
			IsRetain:  m.IsRetain,
			TopicName: m.TopicName,
			Payload:   m.Payload,
		}
	}
}

// SubScribe topic(s)
func (c *client) Subscribe(h SubHandler, topics ...*Topic) {
	for _, t := range topics {
		if v, ok := c.subs[t.Name]; ok {
			// subscribed message before, append this subscriber to the list
			_ = append(v, h)
			c.subs[t.Name] = v
		} else {
			// first time subscribe, start a goroutine to handle msg
		}
	}
}

// UnSubscribe topic(s)
func (c *client) UnSubscribe(h UnSubHandler, topics ...string) {
	for _, t := range topics {
		delete(c.subs, t)
	}
}

// Wait will wait for all connection to exit
// Once called Wait(), you should never add any server to this client
func (c *client) Wait() {
	wg := &sync.WaitGroup{}
	connSet := make([]*connImpl, 0)
	c.conn.Range(func(k, v interface{}) bool {
		connSet = append(connSet, v.(*connImpl))
		return true
	})

	for _, v := range connSet {
		tmp := v
		wg.Add(1)
		go func() {
			tmp.wait()
			wg.Done()
		}()
	}
	lg.d("client waiting")
	wg.Wait()
}

// Destroy will disconnect form all server
// If force is true, then close connection without sending a DisConnPacket
func (c *client) Destroy(force bool) {
	c.conn.Range(func(k, v interface{}) bool {
		va := v.(*connImpl)
		va.close(force)
		return true
	})
}

// connect to one server and start mqtt logic
func (c *client) connect(server string, h ConHandler) {
	var conn net.Conn
	var err error

	if c.options.tlsConfig != nil {
		// with tls
		conn, err = tls.DialWithDialer(&net.Dialer{Timeout: c.options.dialTimeout}, "tcp", server, c.options.tlsConfig)
		if err != nil {
			lg.e("connection with tls failed", err)
			h(server, ConnDialErr)
			return
		}
	} else {
		// without tls
		conn, err = net.DialTimeout("tcp", server, c.options.dialTimeout)
		if err != nil {
			lg.e("connection failed", err)
			h(server, ConnDialErr)
			return
		}
	}

	srvConn := c.newConn(server, conn)
	srvConn.send(context.TODO(), &ConPacket{
		Username:     c.options.username,
		Password:     c.options.password,
		ClientId:     c.options.clientId,
		CleanSession: c.options.cleanSession,
		IsWill:       c.options.isWill,
		WillQos:      c.options.willQos,
		WillTopic:    c.options.willTopic,
		WillMessage:  c.options.willPayload,
		WillRetain:   c.options.willRetain,
		Keepalive:    uint16(c.options.keepalive / time.Second),
	})

	select {
	case pkt, more := <-srvConn.recvC:
		if more {
			if pkt.Type() == CtrlConnAck {
				p := pkt.(*ConAckPacket)
				if p.Code != ConnAccepted {
					h(server, p.Code)
					return
				}
			} else {
				h(server, ConnBadPacket)
				return
			}
		} else {
			h(server, ConnBadPacket)
			return
		}
	case <-time.After(c.options.dialTimeout):
		h(server, ConnTimeout)
		return
	}

	// login success
	// start mqtt logic
	c.conn.Store(server, srvConn)
	srvConn.start()
}

// newConn create a new connection to server
func (c *client) newConn(server string, conn net.Conn) *connImpl {
	co := &connImpl{
		parent:     c,
		name:       server,
		conn:       conn,
		sendBuf:    &bytes.Buffer{},
		keepaliveC: make(chan interface{}),
		recvC:      make(chan Packet),
		sendC:      make(chan sendPkt),
		workers:    &sync.WaitGroup{},
		packetId:   &sync.Map{},
		currentId:  &atomic.Value{},
	}
	co.currentId.Store(uint16(0))

	// handle packet receive
	go func() {
		for {
			pkt, err := decodeOnePacket(conn)
			if err != nil {
				lg.e("connection to server broken: ", err)
				close(co.recvC)
				close(co.sendC)
				close(co.keepaliveC)
				break
			}

			// pass packets
			if pkt == PingRespPacket {
				lg.d("received keepalive message")
				co.keepaliveC <- nil
			} else {
				lg.d("received message")
				co.recvC <- pkt
			}
		}
	}()

	// handle packet send
	go func() {
		for pkt := range co.sendC {
			pkt.packet.Bytes(co.sendBuf)
			if _, err := co.sendBuf.WriteTo(conn); err != nil {
				pkt.errC <- err
				break
			}

			if pkt.packet.Type() == CtrlDisConn {
				// disconnect to server
				break
			}
		}
	}()

	return co
}

// connImpl is the wrapper of connection to server
type connImpl struct {
	parent     *client          // client which created this connection
	name       string           // server addr info
	conn       net.Conn         // connection to server
	sendBuf    *bytes.Buffer    // buffer for packet send
	recvC      chan Packet      // received packet from server
	sendC      chan sendPkt     // need to send packet
	keepaliveC chan interface{} // keepalive packet
	currentId  *atomic.Value
	packetId   *sync.Map       // used pktId (key: packetId, value: packet)
	workers    *sync.WaitGroup // mqtt logic processor
}

type sendPkt struct {
	packet Packet
	errC   chan error
}

// send mqtt packet
func (c *connImpl) send(ctx context.Context, pkt Packet) <-chan error {
	errC := make(chan error)
	c.sendC <- sendPkt{
		packet: pkt,
		errC:   errC,
	}
	return errC
}

// start mqtt logic
func (c *connImpl) start() {
	// start keepalive if required
	if c.parent.options.keepalive > 0 {
		c.workers.Add(1)
		go func() {
			c.keepalive()
			c.workers.Done()
		}()
	}

	// inspect incoming packet
	for pkt := range c.recvC {
		switch pkt.Type() {
		case CtrlPubAck:
		case CtrlPubRecv:
		case CtrlPubRel:
		case CtrlPubComp:
		}
	}
	// TODO: complete mqtt logic
}

// keepalive with server
func (c *connImpl) keepalive() {
	lg.d("started keepalive")
	t := time.NewTicker(c.parent.options.keepalive)
	defer t.Stop()
	defer lg.d("keepalive stopped")

	for range t.C {
		c.send(context.TODO(), PingReqPacket)

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

// get next valid packet id
func (c *connImpl) nextPacketId() uint16 {
	nextId := c.currentId.Load().(uint16) + 1
	for _, ok := c.packetId.Load(nextId); ok; {
		c.currentId.Store(nextId + 1)
	}
	c.currentId.Store(nextId)
	return nextId
}

// close this connection
func (c *connImpl) close(force bool) {
	lg.v(c.name, "close(", force, ")")
	c.conn.Close()
}

// wait for connection lost or close
func (c *connImpl) wait() {
	lg.v(c.name, "wait()")
	c.workers.Wait()
}
