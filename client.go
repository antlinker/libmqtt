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

type Client interface {
}

type client struct {
	username string
	password string
}

func NewClient(options ...Option) Client {
	c := &client{}
	for _, o := range options {
		o(c)
	}

	return c
}

func (c *client) Connect() {

}

func (c *client) Publish(payload []byte) {

}

func (c *client) Subscribe() <-chan Packet {
	return nil
}
