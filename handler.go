package libmqtt

// ConHandler handler the bad connect result
type ConHandler func(server string, code ConAckCode)

// PubHandler handler bad topic pub
type PubHandler func(topic string, code PubAckCode)

// SubHandler handler bad topic sub
type SubHandler func(topic string, qos QosLevel, msg []byte)

// UnSubHandler handler bad topic unSub
type UnSubHandler func(topic string, code SubAckCode)
