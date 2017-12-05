package libmqtt

import (
	"bytes"
)

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
