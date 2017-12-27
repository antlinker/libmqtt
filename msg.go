package libmqtt

type msgType uint8

const (
	pubMsg msgType = iota
	subMsg
	unSubMsg
	netMsg
)

type message struct {
	what msgType
	code byte
	msg  string
	err  error
	obj  interface{}
}

func newPubMsg(topic string, err error) *message {
	return &message{
		what: pubMsg,
		msg:  topic,
		err:  err,
	}
}

func newSubMsg(p []*Topic, err error) *message {
	return &message{
		what: subMsg,
		obj:  p,
		err:  err,
	}
}

func newUnSubMsg(topics []string, err error) *message {
	return &message{
		what: unSubMsg,
		err:  err,
		obj:  topics,
	}
}

func newNetMsg(server string, err error) *message {
	return &message{
		what: netMsg,
		msg:  server,
		err:  err,
	}
}
