package libmqtt

import (
	"io"
)

func decodeRemainLength(buffer io.ByteReader) (result int, err error) {
	m := 1
	var encodedByte byte
	for (encodedByte & 128) != 0 {
		encodedByte, err = buffer.ReadByte()
		if err != nil {
			return
		}
		result += int(encodedByte&127) * m
		m *= 128

		if m > 128*128*128 {
			return
		}
	}

	return
}
