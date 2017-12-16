package libmqtt

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"sync"
	"sync/atomic"
	"time"
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

// NewClient will create a new mqtt client
func NewClient(options ...Option) Client {
	c := defaultClient()
	for _, o := range options {
		o(c)
	}

	return c
}

type clientOptions struct {
	servers         []string      // server address strings
	dialTimeout     time.Duration // dial timeout in second
	clientId        string        // used by ConPacket
	username        string        // used by ConPacket
	password        string        // used by ConPacket
	keepalive       time.Duration // used by ConPacket (time in second)
	keepaliveFactor float64       // used for reasonable amount time to close conn if no ping resp
	cleanSession    bool          // used by ConPacket
	isWill          bool          // used by ConPacket
	willTopic       string        // used by ConPacket
	willPayload     []byte        // used by ConPacket
	willQos         byte          // used by ConPacket
	willRetain      bool          // used by ConPacket
	tlsConfig       *tls.Config
	bf              *BackoffOption
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

	// Close all client connection
	Close()
}

type client struct {
	options     clientOptions
	subC        chan PublishPacket // sub channel
	pubC        chan PublishPacket // pub channel
	subscribers *sync.Map          // topic -> []subscriber
	serverConn  *sync.Map          // server connections
	workers     *sync.WaitGroup    // mqtt logic processor
	packetId    *atomic.Value      // current packet id
}

func defaultClient() *client {
	c := &client{
		options: clientOptions{
			bf: &BackoffOption{
				MaxDelay:   120, // default max retry delay is 2min
				FirstDelay: 1,   // first retry delay is 1s
				Factor:     1.5,
			},
			dialTimeout:     20 * time.Second,
			keepalive:       2 * time.Minute, // default keepalive interval is 2min
			keepaliveFactor: 1.5,             // default reasonable amount of time 3min
		},
		pubC:        make(chan PublishPacket, 256),
		subC:        make(chan PublishPacket, 256),
		subscribers: &sync.Map{},
		serverConn:  &sync.Map{},
		workers:     &sync.WaitGroup{},
		packetId:    &atomic.Value{},
	}
	c.packetId.Store(uint16(0))
	return c
}

func (c *client) Connect(h ConHandler) {
	for _, s := range c.options.servers {
		c.workers.Add(1)
		go c.connect(s, h)
	}
}

func (c *client) Publish(h PubHandler, msg ...*TopicMsg) {
	for _, m := range msg {
		if m.Qos > Qos2 {
			m.Qos = Qos2
		}

		c.pubC <- PublishPacket{
			Qos:       m.Qos,
			IsRetain:  m.IsRetain,
			TopicName: m.TopicName,
			Payload:   m.Payload,
			packetId:  c.nextPacketId(),
		}
	}
}

func (c *client) Subscribe(h SubHandler, topics ...*Topic) {
	for _, t := range topics {
		if v, ok := c.subscribers.Load(t.Name); ok {
			// subscribed message before, append this subscriber to the list
			subs := v.([]SubHandler)
			subs = append(subs, h)
			c.subscribers.Store(t, subs)
		} else {
			// first time subscribe, start a goroutine to handle msg
		}
	}
}

func (c *client) UnSubscribe(h UnSubHandler, topics ...string) {
	for _, t := range topics {
		c.subscribers.Delete(t)
	}
}

func (c *client) Wait() {
	c.workers.Wait()
}

func (c *client) Close() {
	c.serverConn.Range(func(k, v interface{}) bool {
		impl := v.(*connImpl)
		impl.close(false)
		return true
	})
}

// connect to one server and start mqtt logic
func (c *client) connect(server string, h ConHandler) {
	defer c.workers.Done()
	var conn net.Conn
	var err error

	if c.options.tlsConfig != nil {
		// connection with tls
		conn, err = tls.DialWithDialer(&net.Dialer{Timeout: c.options.dialTimeout}, "tcp", server, c.options.tlsConfig)
		if err != nil {
			h(server, ConnDialErr)
			return
		}
	} else {
		// connection without tls
		conn, err = net.DialTimeout("tcp", server, c.options.dialTimeout)
		if err != nil {
			h(server, ConnDialErr)
			return
		}
	}

	srvConn := c.newConn(server, conn)
	srvConn.send(&ConPacket{
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

	// register this conn to client
	c.serverConn.Store(server, srvConn)

	// start mqtt logic
	srvConn.startLogic()
}

// get next valid packet id
func (c *client) nextPacketId() uint16 {
	currentId := c.packetId.Load().(uint16)
	defer c.packetId.Store(currentId + 1)
	return currentId
}

func (c *client) newConn(server string, conn net.Conn) *connImpl {
	co := &connImpl{
		parent:  c,
		name:    server,
		conn:    conn,
		sendBuf: &bytes.Buffer{},
		recvC:   make(chan Packet, 256),
		sendC:   make(chan Packet, 256),
	}

	// handle packet receive
	go func() {
		for {
			pkt, err := decodeOnePacket(conn)
			if err != nil {
				close(co.recvC)
				return
			}
			co.recvC <- pkt
		}
	}()

	// handle packet send
	go func() {
		for pkt := range co.sendC {
			pkt.Bytes(co.sendBuf)
			_, err := co.sendBuf.WriteTo(conn)
			if err != nil {
				return
			}
		}
	}()

	return co
}

type connImpl struct {
	parent         *client
	name           string
	conn           net.Conn
	sendBuf        *bytes.Buffer
	recvC          chan Packet
	sendC          chan Packet
	keepaliveRecvC chan Packet
}

// send mqtt packet
func (c *connImpl) send(pkt Packet) {
	c.sendC <- pkt
}

// start mqtt logic
func (c *connImpl) startLogic() {
	// start keepalive if required
	if c.parent.options.keepalive > 0 {
		go func() {
			t := time.NewTicker(c.parent.options.keepalive)
			for range t.C {
				c.send(PingReqPacket)
				select {
				case <-c.keepaliveRecvC:
					continue
				case <-time.After(c.parent.options.keepalive * time.Duration(c.parent.options.keepaliveFactor)):
					t.Stop()
					c.conn.Close()
				}
			}
		}()
	}

	// inspect incoming packet
	//for pkt := range c.recvC {
	//
	//}
	// TODO: complete mqtt logic
}

// close this connection
func (c *connImpl) close(force bool) {
	// TODO: close up connection
}
