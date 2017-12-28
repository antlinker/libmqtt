package main

import (
	mq "github.com/goiiot/libmqtt"
)

func connHandler(server string, code mq.ConnAckCode, err error) {
	if err != nil {
		println("\nconnect to server error:", err)
	} else if code != mq.ConnAccepted {
		println("\nconnection rejected by server, code:", code)
	} else {
		println("\nconnected to server")
	}
	print(lineStart)
}

func pubHandler(topic string, err error) {
	if err != nil {
		println("\npub", topic, "failed, error =", err)
	} else {
		println("\npub", topic, "success")
	}
	print(lineStart)
}

func subHandler(topics []*mq.Topic, err error) {
	if err != nil {
		println("\nsub", topics, "failed, error =", err)
	} else {
		println("\nsub", topics, "success")
	}
	print(lineStart)
}

func unSubHandler(topics []string, err error) {
	if err != nil {
		println("\nunsub", topics, "failed, error =", err)
	} else {
		println("\nunsub", topics, "success")
	}
	print(lineStart)
}

func netHandler(server string, err error) {
	println("\nconnection to server, error:", err)
	print(lineStart)
}

func topicHandler(topic string, qos mq.QosLevel, msg []byte) {
	println("\n[MSG] topic:", topic, "msg:", string(msg), "qos:", qos)
	print(lineStart)
}

func invalidQos() {
	println("\nqos level should either be 0, 1 or 2")
	print(lineStart)
}
