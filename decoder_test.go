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
	"bufio"
	"bytes"
	"testing"
)

func TestDecodeRemainLength(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.Write([]byte{0x04})
	length, err := decodeRemainLength(buffer)
	if err != nil || length != 0x04 {
		t.Log(length)
		t.Fail()
	}
	buffer.Reset()
}

func TestDecodeOnePacket(t *testing.T) {
	targetBytes := connWillBytes
	buffer := &bytes.Buffer{}
	_, err := buffer.Write(targetBytes)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	reader := bufio.NewReader(buffer)
	pkt, err := decodeOnePacket(reader)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	buffer.Reset()
	switch pkt.(type) {
	case *ConPacket:
		pkt.Bytes(buffer)
		pktBytes := buffer.Bytes()
		if bytes.Compare(pktBytes, targetBytes) != 0 {
			t.Log(pktBytes)
			t.Fail()
		}
	default:
		t.Fail()
	}
}
