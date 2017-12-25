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

import "bytes"

func encodeDataWithLen(data []byte, buffer *bytes.Buffer) {
	lenStr := len(data)
	buffer.WriteByte(byte(lenStr >> 8))
	buffer.WriteByte(byte(lenStr))
	buffer.Write(data)
}

func encodeRemainLength(length int, buffer *bytes.Buffer) {
	if length <= 0 || length > maxMsgSize {
		return
	}

	var encodedByte byte
	for length > 0 {
		encodedByte = byte(length % 128)
		length = length / 128
		if length > 0 {
			encodedByte |= 128
		}
		buffer.WriteByte(encodedByte)
	}
}
