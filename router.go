package libmqtt

import (
	"regexp"
	"sync"
)

// TopicRouter defines how to route the topic message to handler
type TopicRouter interface {
	// Name is the name of router
	Name() string
	Handle(topic string, h SubHandler)
	Dispatch(p PublishPacket)
}

// RegexRouter use regex to match topic messages
type RegexRouter struct {
	sync.Map
}

// Name is the name of router
func (r *RegexRouter) Name() string {
	return "regex"
}

// Handle will register the topic with handler
func (r *RegexRouter) Handle(topicRegex string, h SubHandler) {
	r.Store(regexp.MustCompile(topicRegex), h)
}

// Dispatch the received packet
func (r *RegexRouter) Dispatch(p PublishPacket) {
	r.Range(func(k, v interface{}) bool {
		if reg := k.(*regexp.Regexp); reg.MatchString(p.TopicName) {
			handler := v.(SubHandler)
			handler(p.TopicName, p.Qos, p.Payload)
		}
		return true
	})
}

// TextRouter uses plain string comparison to dispatch topic message
// this is the default router in client
type TextRouter struct {
	sync.Map
}

// Name is the name of router
func (r *TextRouter) Name() string {
	return "text"
}

// Handle will register the topic with handler
func (r *TextRouter) Handle(topic string, h SubHandler) {
	r.Store(topic, h)
}

// Dispatch the received packet
func (r *TextRouter) Dispatch(p PublishPacket) {
	if h, ok := r.Load(p.TopicName); ok {
		handler := h.(SubHandler)
		handler(p.TopicName, p.Qos, p.Payload)
	}
}

// RestRouter is a HTTP RESTFul URL style router
// TODO
type restRouter struct {
}

// Name is the name of router
func (r *restRouter) Name() string {
	return "rest"
}
