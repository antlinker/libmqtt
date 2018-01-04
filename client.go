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
	"errors"
	"io/ioutil"
	"math"
	"net"
	"sync"
	"time"
)

var (
	// ErrTimeOut connection timeout error
	ErrTimeOut = errors.New("connection timeout ")
)

// BackoffOption defines the parameters for the reconnecting backoff strategy.
type BackoffOption struct {
	// MaxDelay defines the upper bound of backoff delay,
	// which is time in second.
	MaxDelay time.Duration

	// FirstDelay is the time to wait before retrying after the first failure,
	// also time in second.
	FirstDelay time.Duration

	// Factor is applied to the backoff after each retry.
	// e.g. FirstDelay = 1 and Factor = 2, then the SecondDelay is 2, the ThirdDelay is 4s
	Factor float64
}

// Option is client option for connection options
type Option func(*client) error

// WithPersist defines the persist method to be used
func WithPersist(method PersistMethod) Option {
	return func(c *client) error {
		if method != nil {
			c.persist = method
		}
		return nil
	}
}

// WithCleanSession will set clean flag in connect packet
func WithCleanSession(f bool) Option {
	return func(c *client) error {
		c.options.cleanSession = f
		return nil
	}
}

// WithIdentity for username and password
func WithIdentity(username, password string) Option {
	return func(c *client) error {
		c.options.username = username
		c.options.password = password
		return nil
	}
}

// WithKeepalive set the keepalive interval (time in second)
func WithKeepalive(keepalive uint16, factor float64) Option {
	return func(c *client) error {
		c.options.keepalive = time.Duration(keepalive) * time.Second
		if factor > 1 {
			c.options.keepaliveFactor = factor
		} else {
			factor = 1.2
		}
		return nil
	}
}

// WithBackoffStrategy will set reconnect backoff strategy
func WithBackoffStrategy(firstDelay, maxDelay time.Duration, factor float64) Option {
	return func(c *client) error {
		if firstDelay < time.Millisecond {
			firstDelay = time.Millisecond
		}

		if maxDelay < firstDelay {
			maxDelay = firstDelay
		}

		if factor < 1 {
			factor = 1
		}

		c.options.backoffOption = &BackoffOption{
			FirstDelay: firstDelay,
			MaxDelay:   maxDelay,
			Factor:     factor,
		}

		return nil
	}
}

// WithClientID set the client id for connection
func WithClientID(clientID string) Option {
	return func(c *client) error {
		c.options.clientID = clientID
		return nil
	}
}

// WithWill mark this connection as a will teller
func WithWill(topic string, qos QosLevel, retain bool, payload []byte) Option {
	return func(c *client) error {
		c.options.isWill = true
		c.options.willTopic = topic
		c.options.willQos = qos
		c.options.willRetain = retain
		c.options.willPayload = payload
		return nil
	}
}

// WithServer adds servers as client server
// Just use "ip:port" or "domain.name:port"
// Only TCP connection supported for now
func WithServer(servers ...string) Option {
	return func(c *client) error {
		c.options.servers = servers
		return nil
	}
}

// WithTLS for client tls certification
func WithTLS(certFile, keyFile string, caCert string, serverNameOverride string, skipVerify bool) Option {
	return func(c *client) error {
		b, err := ioutil.ReadFile(caCert)
		if err != nil {
			return err
		}
		cp := x509.NewCertPool()
		if !cp.AppendCertsFromPEM(b) {
			return err
		}
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return err
		}

		c.options.tlsConfig = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: skipVerify,
			ClientCAs:          cp,
			ServerName:         serverNameOverride,
		}
		return nil
	}
}

// WithDialTimeout for connection time out (time in second)
func WithDialTimeout(timeout uint16) Option {
	return func(c *client) error {
		c.options.dialTimeout = time.Duration(timeout) * time.Second
		return nil
	}
}

// WithSendBuf designate the channel size of send
func WithSendBuf(size int) Option {
	return func(c *client) error {
		if size < 1 {
			size = 1
		} else if size > 1024 {
			size = 1024
		}
		c.options.sendChanSize = size
		return nil
	}
}

// WithRecvBuf designate the channel size of receive
func WithRecvBuf(size int) Option {
	return func(c *client) error {
		if size < 1 {
			size = 1
		} else if size > 1024 {
			size = 1024
		}
		c.options.recvChanSize = size
		return nil
	}
}

// WithRouter set the router for topic dispatch
func WithRouter(r TopicRouter) Option {
	return func(c *client) error {
		if r != nil {
			c.router = r
		}
		return nil
	}
}

// WithLog will create basic logger for log
func WithLog(l LogLevel) Option {
	return func(c *client) error {
		lg = newLogger(l)
		return nil
	}
}

// NewClient will create a new mqtt client
func NewClient(options ...Option) (Client, error) {
	c := defaultClient()

	for _, o := range options {
		err := o(c)
		if err != nil {
			return nil, err
		}
	}

	if len(c.options.servers) < 1 {
		return nil, errors.New("no server provided, won't work ")
	}

	c.msgC = make(chan *message)
	c.sendC = make(chan *bytes.Buffer, c.options.sendChanSize)
	c.recvC = make(chan *PublishPacket, c.options.recvChanSize)

	return c, nil
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
	backoffOption   *BackoffOption // backoff option for client reconnection
}

// Client act as a mqtt client
type Client interface {
	// Handle register topic handlers, mostly used for RegexHandler, RestHandler
	// the default handler inside the client is TextHandler, which match the exactly same topic
	Handle(topic string, h TopicHandler)

	// Connect to all specified server with client options
	Connect(ConnHandler)

	// Publish a message for the topic
	Publish(packets ...*PublishPacket)

	// Subscribe topic(s)
	Subscribe(topics ...*Topic)

	// UnSubscribe topic(s)
	UnSubscribe(topics ...string)

	// Wait will wait until all connection finished
	Wait()

	// Destroy all client connection
	Destroy(force bool)

	// handlers
	HandlePub(PubHandler)
	HandleSub(SubHandler)
	HandleUnSub(UnSubHandler)
	HandleNet(NetHandler)
	HandlePersist(PersistHandler)
}

type client struct {
	options *clientOptions      // client connection options
	subs    *sync.Map           // Topic(s) -> []TopicHandler
	conn    *sync.Map           // ServerAddr -> connection
	bufPool *sync.Pool          // a pool for buffer
	msgC    chan *message       // error channel
	sendC   chan *bytes.Buffer  // Pub channel for sending publish packet to server
	recvC   chan *PublishPacket // recv channel for server pub receiving
	idGen   *idGenerator        // Packet id generator
	router  TopicRouter         // Topic router
	persist PersistMethod       // Persist method
	workers *sync.WaitGroup     // Workers (connections)

	// success/error handlers
	pH  PubHandler
	sH  SubHandler
	uH  UnSubHandler
	nH  NetHandler
	psH PersistHandler
}

// defaultClient create the client with default options
func defaultClient() *client {
	return &client{
		options: &clientOptions{
			sendChanSize: 128,
			recvChanSize: 128,
			backoffOption: &BackoffOption{
				MaxDelay:   2 * time.Minute, // default max retry delay is 2min
				FirstDelay: 5 * time.Second, // first retry delay is 5s
				Factor:     1.5,
			},
			dialTimeout:     20 * time.Second, // default timeout when dial to server
			keepalive:       2 * time.Minute,  // default keepalive interval is 2min
			keepaliveFactor: 1.5,              // default reasonable amount of time 3min
		},
		router:  NewTextRouter(),
		subs:    &sync.Map{},
		conn:    &sync.Map{},
		idGen:   newIDGenerator(),
		workers: &sync.WaitGroup{},
		persist: &NonePersist{},
		bufPool: &sync.Pool{New: func() interface{} { return &bytes.Buffer{} }},
	}
}

// Handle subscription message route
func (c *client) Handle(topic string, h TopicHandler) {
	if h != nil {
		lg.d("HANDLE registered handler, topic =", topic)
		c.router.Handle(topic, h)
	}
}

// Connect to all designated server
func (c *client) Connect(h ConnHandler) {
	lg.d("CLIENT connect to server, handler =", h)
	go func() {
		for pkt := range c.recvC {
			c.router.Dispatch(pkt)
		}
	}()

	go func() {
		for m := range c.msgC {
			switch m.what {
			case pubMsg:
				if c.pH != nil {
					go c.pH(m.msg, m.err)
				}
			case subMsg:
				if c.sH != nil {
					go c.sH(m.obj.([]*Topic), m.err)
				}
			case unSubMsg:
				if c.uH != nil {
					go c.uH(m.obj.([]string), m.err)
				}
			case netMsg:
				if c.nH != nil {
					go c.nH(m.msg, m.err)
				}
			case persistMsg:
				if c.psH != nil {
					go c.psH(m.err)
				}
			}
		}
	}()

	for _, s := range c.options.servers {
		c.workers.Add(1)
		go c.connect(s, h, 0)
	}
}

// Publish message(s) to topic(s), one to one
func (c *client) Publish(msg ...*PublishPacket) {
	buf := c.getBuf()
	for _, m := range msg {
		if m == nil {
			continue
		}

		p := m
		if p.Qos > Qos2 {
			p.Qos = Qos2
		}

		if p.Qos != Qos0 {
			p.PacketID = c.idGen.next(p)
			c.persist.Store(sendKey(p.PacketID), p)
		} else {
			c.msgC <- newPubMsg(p.TopicName, nil)
		}
		p.Bytes(buf)
	}
	c.sendC <- buf
}

// SubScribe topic(s)
func (c *client) Subscribe(topics ...*Topic) {
	lg.d("CLIENT subscribe, topic(s) =", topics)
	s := &SubscribePacket{Topics: topics}
	s.PacketID = c.idGen.next(s)
	buf := c.getBuf()
	s.Bytes(buf)
	c.sendC <- buf
	c.persist.Store(sendKey(s.PacketID), s)
}

// UnSubscribe topic(s)
func (c *client) UnSubscribe(topics ...string) {
	lg.d("CLIENT unsubscribe topic(s) =", topics)
	for _, t := range topics {
		c.subs.Delete(t)
	}
	u := &UnSubPacket{
		TopicNames: topics,
	}
	u.PacketID = c.idGen.next(u)
	buf := c.getBuf()
	u.Bytes(buf)
	c.sendC <- buf
	c.persist.Store(sendKey(u.PacketID), u)
}

// Wait will wait for all connection to exit
func (c *client) Wait() {
	lg.i("CLIENT wait for all connections")
	c.workers.Wait()
}

// Destroy will disconnect form all server
// If force is true, then close connection without sending a DisConnPacket
func (c *client) Destroy(force bool) {
	lg.d("CLIENT destroying client with force =", force)
	// TODO close all channel properly
	c.options.backoffOption = nil
	if force {
		c.conn.Range(func(k, v interface{}) bool {
			va := v.(*connImpl)
			va.conn.Close()
			return true
		})
	} else {
		c.conn.Range(func(k, v interface{}) bool {
			va := v.(*connImpl)
			va.close()
			return true
		})
	}
}

// HandlePubMsg register handler for pub error
func (c *client) HandlePub(h PubHandler) {
	lg.d("CLIENT registered pub handler")
	c.pH = h
}

// HandleSubMsg register handler for extra sub info
func (c *client) HandleSub(h SubHandler) {
	lg.d("CLIENT registered sub handler")
	c.sH = h
}

// HandleUnSubMsg register handler for unsubscription error
func (c *client) HandleUnSub(h UnSubHandler) {
	lg.d("CLIENT registered unsub handler")
	c.uH = h
}

// HandleNet register handler for net error
func (c *client) HandleNet(h NetHandler) {
	lg.d("CLIENT registered net handler")
	c.nH = h
}

// HandleNet register handler for net error
func (c *client) HandlePersist(h PersistHandler) {
	lg.d("CLIENT registered persist handler")
	c.psH = h
}

// connect to one server and start mqtt logic
func (c *client) connect(server string, h ConnHandler, reconnectDelay time.Duration) {
	defer c.workers.Done()
	var conn net.Conn
	var err error

	if c.options.tlsConfig != nil {
		// with tls
		conn, err = tls.DialWithDialer(&net.Dialer{Timeout: c.options.dialTimeout}, "tcp", server, c.options.tlsConfig)
		if err != nil {
			lg.e("CLIENT connect with tls failed, err =", err, "server =", server)
			if h != nil {
				h(server, math.MaxUint8, err)
			}
			return
		}
	} else {
		// without tls
		conn, err = net.DialTimeout("tcp", server, c.options.dialTimeout)
		if err != nil {
			lg.e("CLIENT connect failed, err =", err, "server =", server)
			if h != nil {
				h(server, math.MaxUint8, err)
			}
			return
		}
	}

	connImpl := &connImpl{
		parent:     c,
		name:       server,
		conn:       conn,
		clientBuf:  &bytes.Buffer{},
		sendBuf:    &bytes.Buffer{},
		keepaliveC: make(chan int),
		logicSendC: make(chan Packet),
		netRecvC:   make(chan Packet),
	}

	go connImpl.handleLogicSend()
	go connImpl.handleRecv()
	go connImpl.handleClientSend()

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
	case pkt, more := <-connImpl.netRecvC:
		if more {
			if pkt.Type() == CtrlConnAck {
				p := pkt.(*ConAckPacket)
				if p.Code != ConnAccepted {
					if h != nil {
						h(server, p.Code, nil)
					}
					return
				}
			} else {
				if h != nil {
					h(server, math.MaxUint8, ErrBadPacket)
				}
				return
			}
		} else {
			if h != nil {
				h(server, math.MaxUint8, ErrBadPacket)
			}
			return
		}
	case <-time.After(c.options.dialTimeout):
		if h != nil {
			h(server, math.MaxUint8, ErrTimeOut)
		}
		return
	}

	lg.i("CLIENT connected server =", server)
	if h != nil {
		go h(server, ConnAccepted, nil)
	}

	// login success, start mqtt logic
	c.conn.Store(server, connImpl)
	connImpl.logic()

	if c.options.backoffOption != nil {
		c.workers.Add(1)
		lg.w("CLIENT reconnecting to server =", server, "delay =", reconnectDelay)
		go func() {
			time.Sleep(reconnectDelay)
			reconnectDelay = time.Duration(float64(reconnectDelay) * c.options.backoffOption.Factor)
			if reconnectDelay > c.options.backoffOption.MaxDelay {
				reconnectDelay = c.options.backoffOption.MaxDelay
			}
			c.connect(server, h, reconnectDelay)
		}()
	}
}

// get a buffer from buf pool
func (c *client) getBuf() *bytes.Buffer {
	return c.bufPool.Get().(*bytes.Buffer)
}

// put back the buffer
func (c *client) putBuf(buf *bytes.Buffer) {
	if buf == nil {
		return
	}
	c.bufPool.Put(buf)
}

// connImpl is the wrapper of connection to server
// tend to actual packet send and receive
type connImpl struct {
	parent     *client       // client which created this connection
	name       string        // server addr info
	conn       net.Conn      // connection to server
	sendBuf    *bytes.Buffer // buffer for logic packet send
	clientBuf  *bytes.Buffer // buffer for client packet send
	logicSendC chan Packet   // logic send channel
	netRecvC   chan Packet   // received packet from server
	keepaliveC chan int      // keepalive packet
}

// start mqtt logic
func (c *connImpl) logic() {
	// start keepalive if required
	if c.parent.options.keepalive > 0 {
		go c.keepalive()
	}

	// inspect incoming packet
	for pkt := range c.netRecvC {
		switch pkt.Type() {
		case CtrlSubAck:
			p := pkt.(*SubAckPacket)
			lg.d("NET received SubAck, id =", p.PacketID)

			if originPkt, ok := c.parent.idGen.getExtra(p.PacketID); ok {
				switch originPkt.(type) {
				case *SubscribePacket:
					originSub := originPkt.(*SubscribePacket)
					N := len(p.Codes)
					for i, v := range originSub.Topics {
						if i < N {
							v.Qos = p.Codes[i]
						}
					}
					c.parent.msgC <- newSubMsg(originSub.Topics, nil)
					c.parent.idGen.free(p.PacketID)

					c.parent.persist.Delete(sendKey(p.PacketID))
				}
			}
		case CtrlUnSubAck:
			p := pkt.(*UnSubAckPacket)
			lg.d("NET received UnSubAck, id =", p.PacketID)

			if originPkt, ok := c.parent.idGen.getExtra(p.PacketID); ok {
				switch originPkt.(type) {
				case *UnSubPacket:
					originUnSub := originPkt.(*UnSubPacket)
					c.parent.msgC <- newUnSubMsg(originUnSub.TopicNames, nil)
					c.parent.idGen.free(p.PacketID)

					c.parent.persist.Delete(sendKey(p.PacketID))
				}
			}
		case CtrlPublish:
			p := pkt.(*PublishPacket)
			lg.d("NET received publish, id =", p.PacketID, "QoS =", p.Qos)
			// received server publish, send to client
			c.parent.recvC <- p

			// tend to QoS
			switch p.Qos {
			case Qos1:
				lg.d("NET send PubAck for Publish, id =", p.PacketID)
				c.send(&PubAckPacket{PacketID: p.PacketID})

				c.parent.persist.Store(recvKey(p.PacketID), pkt)
			case Qos2:
				lg.d("NET send PubRecv for Publish, id =", p.PacketID)
				c.send(&PubRecvPacket{PacketID: p.PacketID})

				c.parent.persist.Store(recvKey(p.PacketID), pkt)
			}
		case CtrlPubAck:
			p := pkt.(*PubAckPacket)
			lg.d("NET received PubAck, id =", p.PacketID)

			if originPkt, ok := c.parent.idGen.getExtra(p.PacketID); ok {
				switch originPkt.(type) {
				case *PublishPacket:
					originPub := originPkt.(*PublishPacket)
					if originPub.Qos == Qos1 {
						c.parent.msgC <- newPubMsg(originPub.TopicName, nil)
						c.parent.idGen.free(p.PacketID)

						c.parent.persist.Delete(sendKey(p.PacketID))
					}
				}
			}
		case CtrlPubRecv:
			p := pkt.(*PubRecvPacket)
			lg.d("NET received PubRec, id =", p.PacketID)

			if originPkt, ok := c.parent.idGen.getExtra(p.PacketID); ok {
				switch originPkt.(type) {
				case *PublishPacket:
					originPub := originPkt.(*PublishPacket)
					if originPub.Qos == Qos2 {
						c.send(&PubRelPacket{PacketID: p.PacketID})
						lg.d("NET send PubRel, id =", p.PacketID)
					}
				}
			}
		case CtrlPubRel:
			p := pkt.(*PubRelPacket)
			lg.d("NET send PubRel, id =", p.PacketID)

			if originPkt, ok := c.parent.idGen.getExtra(p.PacketID); ok {
				switch originPkt.(type) {
				case *PublishPacket:
					originPub := originPkt.(*PublishPacket)
					if originPub.Qos == Qos2 {
						c.send(&PubCompPacket{PacketID: p.PacketID})
						lg.d("NET send PubComp, id =", p.PacketID)

						c.parent.persist.Store(recvKey(p.PacketID), pkt)
					}
				}
			}
		case CtrlPubComp:
			p := pkt.(*PubCompPacket)
			lg.d("NET received PubComp, id =", p.PacketID)

			if originPkt, ok := c.parent.idGen.getExtra(p.PacketID); ok {
				switch originPkt.(type) {
				case *PublishPacket:
					originPub := originPkt.(*PublishPacket)
					if originPub.Qos == Qos2 {
						c.send(&PubRelPacket{PacketID: p.PacketID})
						lg.d("NET send PubRel, id =", p.PacketID)

						c.parent.msgC <- newPubMsg(originPub.TopicName, nil)
						c.parent.idGen.free(p.PacketID)

						c.parent.persist.Delete(sendKey(p.PacketID))
					}
				}
			}
		default:
			lg.d("NET received packet, type =", pkt.Type())
		}
	}
}

// keepalive with server
func (c *connImpl) keepalive() {
	lg.d("NET start keepalive")

	t := time.NewTicker(c.parent.options.keepalive * 3 / 4)
	timeout := time.Duration(float64(c.parent.options.keepalive) * c.parent.options.keepaliveFactor)
	timeoutTimer := time.NewTimer(timeout)
	defer t.Stop()

	for range t.C {
		c.send(PingReqPacket)

		select {
		case _, more := <-c.keepaliveC:
			if !more {
				return
			}
			timeoutTimer.Reset(timeout)
		case <-timeoutTimer.C:
			lg.i("NET keepalive timeout")
			t.Stop()
			c.conn.Close()
			return
		}
	}

	lg.d("NET stop keepalive")
}

// close this connection
func (c *connImpl) close() {
	lg.i("NET connection to server closed, remote =", c.name)
	c.send(DisConPacket)
}

// handle client message send
func (c *connImpl) handleClientSend() {
	for buf := range c.parent.sendC {
		_, err := buf.WriteTo(c.conn)
		if err != nil {
			// DO NOT NOTIFY net err HERE
			// ALWAYS DETECT net err in receive
			buf.Reset()
			break
		}

		c.parent.putBuf(buf)
	}
}

// handle mqtt logic control packet send
func (c *connImpl) handleLogicSend() {
	for logicPkt := range c.logicSendC {
		logicPkt.Bytes(c.sendBuf)
		if _, err := c.sendBuf.WriteTo(c.conn); err != nil {
			// DO NOT NOTIFY net err HERE
			// ALWAYS DETECT net err in receive
			break
		}

		switch logicPkt.Type() {
		case CtrlPubRel:
			c.parent.persist.Store(sendKey(logicPkt.(*PubRelPacket).PacketID), logicPkt)
		case CtrlPubAck:
			c.parent.persist.Delete(sendKey(logicPkt.(*PubAckPacket).PacketID))
		case CtrlPubComp:
			c.parent.persist.Delete(sendKey(logicPkt.(*PubCompPacket).PacketID))
		case CtrlDisConn:
			// disconnect to server
			lg.i("disconnect to server")
			c.conn.Close()
			break
		}
	}
}

// handle all message receive
func (c *connImpl) handleRecv() {
	for {
		pkt, err := DecodeOnePacket(c.conn)
		if err != nil {
			lg.e("NET connection broken, server =", c.name, "err =", err)
			close(c.netRecvC)
			close(c.keepaliveC)
			// TODO send proper net error to net handler
			if err != ErrBadPacket {
				// c.parent.msgC <- newNetMsg(c.name, err)
			}
			break
		}

		if pkt == PingRespPacket {
			lg.d("NET received keepalive message")
			c.keepaliveC <- 1
		} else {
			c.netRecvC <- pkt
		}
	}
}

// send mqtt logic packet
func (c *connImpl) send(pkt Packet) {
	c.logicSendC <- pkt
}
