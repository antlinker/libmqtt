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

// NewStandardRouter will create a standard mqtt router
func NewStandardRouter() *StandardRouter {
	return &StandardRouter{m: &sync.Map{}}
}

// StandardRouter implements standard MQTT routing behaviour
type StandardRouter struct {
	m *sync.Map
}

// Name is the name of router
func (s *StandardRouter) Name() string {
	if s == nil {
		return "<nil>"
	}
	return "StandardRouter"
}

// Handle defines how to register topic with handler
func (s *StandardRouter) Handle(topic string, h TopicHandler) {

}

// Dispatch defines the action to dispatch published packet
func (s *StandardRouter) Dispatch(p *PublishPacket) {

}

// NewRegexRouter will create a regex router
func NewRegexRouter() *RegexRouter {
	return &RegexRouter{m: &sync.Map{}}
}

// RegexRouter use regex to match topic messages
type RegexRouter struct {
	m *sync.Map
}

// Name is the name of router
func (r *RegexRouter) Name() string {
	if r == nil {
		return "<nil>"
	}

	return "RegexRouter"
}

// Handle will register the topic with handler
func (r *RegexRouter) Handle(topicRegex string, h TopicHandler) {
	if r == nil || r.m == nil {
		return
	}
	r.m.Store(regexp.MustCompile(topicRegex), h)
}

// Dispatch the received packet
func (r *RegexRouter) Dispatch(p *PublishPacket) {
	if r == nil || r.m == nil {
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

// NewTextRouter will create a text based router
func NewTextRouter() *TextRouter {
	return &TextRouter{m: &sync.Map{}}
}

// TextRouter uses plain string comparison to dispatch topic message
// this is the default router in client
type TextRouter struct {
	m *sync.Map
}

// Name of TextRouter is "TextRouter"
func (r *TextRouter) Name() string {
	if r == nil {
		return "<nil>"
	}

	return "TextRouter"
}

// Handle will register the topic with handler
func (r *TextRouter) Handle(topic string, h TopicHandler) {
	if r == nil || r.m == nil {
		return
	}

	r.m.Store(topic, h)
}

// Dispatch the received packet
func (r *TextRouter) Dispatch(p *PublishPacket) {
	if r == nil || r.m == nil {
		return
	}

	if h, ok := r.m.Load(p.TopicName); ok {
		handler := h.(TopicHandler)
		handler(p.TopicName, p.Qos, p.Payload)
	}
}
