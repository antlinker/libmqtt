package main

import (
	"strconv"
	"strings"

	mq "github.com/goiiot/libmqtt"
)

func execPub(args []string) bool {
	if client == nil {
		println("please connect to server first")
		return true
	}

	pubs := make([]*mq.PublishPacket, 0)
	for _, v := range args {
		pubStr := strings.Split(v, ",")
		if len(pubStr) != 3 {
			pubUsage()
			return true
		}
		qos, err := strconv.Atoi(pubStr[1])
		if err != nil {
			pubUsage()
			return false
		}
		pubs = append(pubs, &mq.PublishPacket{
			TopicName: pubStr[0],
			Qos:       mq.QosLevel(qos),
			Payload:   []byte(pubStr[2]),
		})
	}
	client.Publish(pubs...)
	return true
}

func pubUsage() {
	println(`p, pub [topic,qos,message] [...] - publish topic message(s)`)
}
