package libmqtt

type Option func(*client)

func WithIdentity(username, password string) Option {
	return func(c *client) {
		c.username = username
		c.password = password
	}
}

func WithTLS() Option {
	return nil
}

// NewClient will create a new mqtt client
func NewClient(options ...Option) Client {
	c := &client{}
	for _, o := range options {
		o(c)
	}

	return c
}

type Client interface {
}

type client struct {
	username string
	password string
}

func (c *client) Connect() {

}

func (c *client) Publish(topic []Topic, payload []byte) {

}

func (c *client) Subscribe(topics []Topic) <-chan Packet {
	return nil
}
