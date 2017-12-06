package libmqtt

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"sync"
	"time"
	"bytes"
)

// BackoffConfig defines the parameters for the reconnecting backoff strategy.
type BackoffConfig struct {
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

// WithKeepalive
func WithKeepalive(keepalive uint16) Option {
	return func(c *client) {
		c.options.keepalive = keepalive
	}
}

// WithBackoffStrategy will set reconnect backoff strategy
func WithBackoffStrategy(bf *BackoffConfig) Option {
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

// WithWill make this connection as a will teller
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

// WithTLS for ssl certification
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
	servers      []string      // server address strings
	dialTimeout  time.Duration // dial timeout in second
	clientId     string        // used by ConnPacket
	username     string        // used by ConnPacket
	password     string        // used by ConnPacket
	keepalive    uint16        // used by ConnPacket
	cleanSession bool          // used by ConnPacket
	isWill       bool          // used by ConnPacket
	willTopic    string        // used by ConnPacket
	willPayload  []byte        // used by ConnPacket
	willQos      byte          // used by ConnPacket
	willRetain   bool          // used by ConnPacket
	tlsConfig    *tls.Config
	bf           *BackoffConfig
}

// Client act as a mqtt client
type Client interface {
	Connect(errHandler ConnHandler)
	Publish(topic string, qos QosLevel, isRetain bool, payload []byte)
	Wait()
}

type client struct {
	options     clientOptions
	sub         chan PublishPacket // sub channel
	pub         chan PublishPacket // pub channel
	subscribers *sync.Map          // topic -> []subscriber
	serverConn  *sync.Map          // server connections
	workers     *sync.WaitGroup    // mqtt logic processor
}

func defaultClient() *client {
	return &client{
		options: clientOptions{
			bf: &BackoffConfig{
				MaxDelay:   120, // default max retry delay is 2min
				FirstDelay: 1,   // first retry delay is 1s
				Factor:     1.5,
			},
			dialTimeout: 20 * time.Second,
			keepalive:   120, // default keepalive interval is 2min
		},
		pub:         make(chan PublishPacket, 256),
		sub:         make(chan PublishPacket, 256),
		subscribers: &sync.Map{},
		serverConn:  &sync.Map{},
		workers:     &sync.WaitGroup{},
	}
}

type ConnHandler func(server string, code ConnAckCode)

// Connect to all specified server with client options
func (c *client) Connect(errHandler ConnHandler) {
	for _, s := range c.options.servers {
		c.workers.Add(1)
		go c.connect(s, errHandler)
	}
}

// Publish a message for the topic
func (c *client) Publish(topic string, qos QosLevel, isRetain bool, payload []byte) {
	if qos > Qos2 {
		qos = Qos2
	}

	c.pub <- PublishPacket{
		Qos:       qos,
		IsRetain:  isRetain,
		TopicName: topic,
		PacketId:  0,
		Payload:   payload,
	}
}

// Publish subscribe to topic(s)
func (c *client) Subscribe(subscriber Subscriber, topics ...Topic) {
	for _, t := range topics {
		if v, ok := c.subscribers.Load(t.Name); ok {
			// subscribed message before, append this subscriber to the list
			subs := v.([]Subscriber)
			subs = append(subs, subscriber)
			c.subscribers.Store(t, subs)
		} else {
			// first time subscribe, start a goroutine to handle msg
		}
	}
}

// UnSubscribe topic(s)
func (c *client) UnSubscribe(topics ...string) {
	for _, t := range topics {
		c.subscribers.Delete(t)
	}
}

func (c *client) connect(server string, errHandler ConnHandler) {
	defer c.workers.Done()

	var conn net.Conn
	var err error

	if c.options.tlsConfig != nil {
		// connection with tls
		conn, err = tls.DialWithDialer(&net.Dialer{Timeout: c.options.dialTimeout}, "tcp", server, c.options.tlsConfig)
		if err != nil {
			errHandler(server, ConnDialErr)
			return
		}
	} else {
		// connection without tls
		conn, err = net.DialTimeout("tcp", server, c.options.dialTimeout)
		if err != nil {
			errHandler(server, ConnDialErr)
			return
		}
	}

	buf := &bytes.Buffer{}

	// send conn packet
	connPkt := &ConnPacket{
		Username:     c.options.username,
		Password:     c.options.password,
		ClientId:     c.options.clientId,
		CleanSession: c.options.cleanSession,
		IsWill:       c.options.isWill,
		WillQos:      c.options.willQos,
		WillTopic:    c.options.willTopic,
		WillMessage:  c.options.willPayload,
		WillRetain:   c.options.willRetain,
		Keepalive:    c.options.keepalive,
	}
	connPkt.Bytes(buf)
	_, err = buf.WriteTo(conn)
	buf.Reset()
	if err != nil {
		errHandler(server, ConnNetErr)
		return
	}

	// receive ConnAck
	pkt, err := decodeOnePacket(conn)
	if err != nil {
		errHandler(server, ConnBadPacket)
		return
	}

	if pkt.Type() == CtrlConnAck {
		p := pkt.(*ConnAckPacket)
		if p.Code != ConnAccepted {
			errHandler(server, p.Code)
			return
		}
	} else {
		errHandler(server, ConnBadPacket)
		return
	}

	c.startLogic(conn, buf)
}

func (c *client) startLogic(conn net.Conn, buffer *bytes.Buffer) {

}

func (c *client) Wait() {
	c.workers.Wait()
}
