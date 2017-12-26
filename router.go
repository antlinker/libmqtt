package libmqtt

import (
	"regexp"
	"sync"
)

type TopicRouter interface {
	Name() string
	Handle(topic string, h SubHandler)
	Dispatch(p PublishPacket)
}

// RegexRouter use regex to match topic messages
type RegexRouter struct {
	sync.Map
}

func (r *RegexRouter) Name() string {
	return "regex"
}

func (r *RegexRouter) Handle(topicRegex string, h SubHandler) {
	r.Store(regexp.MustCompile(topicRegex), h)
}

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

func (r *TextRouter) Name() string {
	return "text"
}

func (r *TextRouter) Handle(topic string, h SubHandler) {
	r.Store(topic, h)
}

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

func (r *restRouter) Name() string {
	return "rest"
}

func (r *restRouter) Match(topic string) bool {
	return true
}
