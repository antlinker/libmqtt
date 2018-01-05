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
	"testing"
)

func TestSubscribePacket_Bytes(t *testing.T) {
	for i, p := range testSubMsgs {
		testBytes(p, testSubMsgBytes[i], t)
	}
}

func TestSubAckPacket_Bytes(t *testing.T) {
	for i, p := range testSubAckMsgs {
		testBytes(p, testSubAckMsgBytes[i], t)
	}
}

func TestUnSubPacket_Bytes(t *testing.T) {
	for i, p := range testUnSubMsgs {
		testBytes(p, testUnSubMsgBytes[i], t)
	}
}

func TestUnSubAckPacket_Bytes(t *testing.T) {
	testBytes(testUnSubAckMsg, testUnSubAckMsgBytes, t)
}
