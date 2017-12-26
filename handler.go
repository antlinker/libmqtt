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

// PubHandler handler bad topic pub
// topic is the client user provided topic
// err is the error happened when publish
type PubHandler func(topic string, err error)

// SubHandler handler topic sub
// topic is the client user provided topic
// code can be SubOkMaxQos0, SubOkMaxQos1, SubOkMaxQos2, SubFail
type SubHandler func(topic string, code SubAckCode, msg []byte)

// UnSubHandler handler topic unSub
// topic is the client user provided topic
// err is the error happened when unsubscribe
type UnSubHandler func(topic string, err error)
