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

package extension

import (
	lib "github.com/goiiot/libmqtt"
)

// HttpRouter is a HTTP URL style router
type HttpRouter struct {
}

// Name of HttpRouter is "HttpRouter"
func (r *HttpRouter) Name() string {
	if r == nil {
		return "<nil>"
	}

	return "HttpRouter"
}

// Handle the topic with TopicHandler h
func (r *HttpRouter) Handle(topic string, h lib.TopicHandler) {

}

// Dispatch the received packet
func (r *HttpRouter) Dispatch(p *lib.PublishPacket) {

}
