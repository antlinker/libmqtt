package main

import (
	"strconv"
	"strings"

	mq "github.com/goiiot/libmqtt"
)

func execSub(args []string) bool {
	if client == nil {
		println("please connect to server first")
		return true
	}

	topics := make([]*mq.Topic, 0)
	for _, v := range args {
		topicStr := strings.Split(v, ",")
		if len(topicStr) != 2 {
			subUsage()
			return true
		}
		qos, err := strconv.Atoi(topicStr[1])
		if err != nil {
			subUsage()
			return true
		}
		topics = append(topics, &mq.Topic{Name: topicStr[0], Qos: mq.QosLevel(qos)})
	}
	for _, t := range topics {
		client.Handle(t.Name, topicHandler)
	}
	client.Subscribe(topics...)
	return true
}

func execUnSub(args []string) bool {
	if client == nil {
		println("please connect to server first")
		return true
	}

	client.UnSubscribe(args...)
	return true
}

func subUsage() {
	println(`s, sub [topic,qos] [...] - subscribe topic(s)`)
}

func unSubUsage() {
	println(`u, unsub [topic] [...] - unsubscribe topic(s)`)
}
