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

type msgType uint8

const (
	pubMsg msgType = iota
	subMsg
	unSubMsg
	netMsg
	persistMsg
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

func newPersistMsg(err error) *message {
	return &message{
		what: persistMsg,
		err:  err,
	}
}
