package libmqtt

import (
	"regexp"
	"sync"
)

// TopicRouter defines how to route the topic message to handler
type TopicRouter interface {
	// Name is the name of router
	Name() string
	// Handle defines how to register topic with handler
	Handle(topic string, h TopicHandler)
	// Dispatch defines the action to dispatch published packet
	Dispatch(p *PublishPacket)
}

// RegexRouter use regex to match topic messages
type RegexRouter struct {
	m *sync.Map
}

// Name is the name of router
func (r *RegexRouter) Name() string {
	return "regex"
}

// Handle will register the topic with handler
func (r *RegexRouter) Handle(topicRegex string, h TopicHandler) {
	if r.m == nil {
		r.m = &sync.Map{}
	}

	r.m.Store(regexp.MustCompile(topicRegex), h)
}

// Dispatch the received packet
func (r *RegexRouter) Dispatch(p *PublishPacket) {
	if r.m == nil {
		return
	}

	r.m.Range(func(k, v interface{}) bool {
		if reg := k.(*regexp.Regexp); reg.MatchString(p.TopicName) {
			handler := v.(TopicHandler)
			handler(p.TopicName, p.Qos, p.Payload)
		}
		return true
	})
}

// TextRouter uses plain string comparison to dispatch topic message
// this is the default router in client
type TextRouter struct {
	m *sync.Map
}

// Name is the name of router
func (r *TextRouter) Name() string {
	return "text"
}

// Handle will register the topic with handler
func (r *TextRouter) Handle(topic string, h TopicHandler) {
	if r.m == nil {
		r.m = &sync.Map{}
	}

	r.m.Store(topic, h)
}

// Dispatch the received packet
func (r *TextRouter) Dispatch(p *PublishPacket) {
	if r.m == nil {
		return
	}

	if h, ok := r.m.Load(p.TopicName); ok {
		handler := h.(TopicHandler)
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

// Handle will register the topic with handler
// TODO
func (r *restRouter) Handle(topic string, h TopicHandler) {

}

// Dispatch the received packet
// TODO
func (r *restRouter) Dispatch(p *PublishPacket) {

}
