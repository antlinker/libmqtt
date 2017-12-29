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

// ConnHandler is the handler which tend to the Connect result
// server is the server address provided by user in client creation call
// code is the ConnResult code
// err is the error happened when connect to server, if a error happened,
// the code value will max byte value (255)
type ConnHandler func(server string, code ConnAckCode, err error)

// TopicHandler handles topic sub message
// topic is the client user provided topic
// code can be SubOkMaxQos0, SubOkMaxQos1, SubOkMaxQos2, SubFail
type TopicHandler func(topic string, qos QosLevel, msg []byte)

// PubHandler handles the error occurred when publish some message
// if err is not nil, that means a error occurred when sending pub msg
type PubHandler func(topic string, err error)

// SubHandler handles the error occurred when subscribe some topic
// if err is not nil, that means a error occurred when sending sub msg
type SubHandler func(topics []*Topic, err error)

// UnSubHandler handles the error occurred when publish some message
type UnSubHandler func(topic []string, err error)

// NetHandler handles the error occurred when net broken
type NetHandler func(server string, err error)

// PersistHandler handles err happened when persist process has trouble
type PersistHandler func(err error)
