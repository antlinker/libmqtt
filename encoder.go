package libmqtt

import (
	"bytes"
)

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
