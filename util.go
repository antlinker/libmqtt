package libmqtt

func boolToByte(flag bool) (result byte) {
	if flag {
		result = 0x01
	}
	return
}
