package cmd

import (
	"github.com/goiiot/libmqtt"
)

func Example() {
	client := libmqtt.NewClient(
		// server address(es)
		libmqtt.WithServer("localhost:1883"),
	)

	// define your topic handlers like a golang http server
	client.Handle("foo", func(topic string, qos libmqtt.QosLevel, msg []byte) {
		// handle the foo topic message
	})

	client.Handle("bar", func(topic string, qos libmqtt.QosLevel, msg []byte) {
		// handle the bar topic message
	})

	// connect to server
	client.Connect(func(server string, code libmqtt.ConnAckCode, err error) {
		if err != nil {
			// failed
			panic(err)
		}

		if code != libmqtt.ConnAccepted {
			// server rejected or in error
			panic(code)
		}

		// success
		// you are now connected to the `server`
		// start your business logic here or send a signal to your logic to start

		// subscribe some topic(s)
		client.Subscribe(
			&libmqtt.Topic{Name: "foo"},
			&libmqtt.Topic{Name: "bar", Qos: libmqtt.Qos1},
			// ...
		)

		// publish some topic message(s)
		client.Publish(
			&libmqtt.PublishPacket{
				TopicName: "foo",
				Qos:       libmqtt.Qos0,
				Payload:   []byte("foo data"),
			},
			&libmqtt.PublishPacket{
				TopicName: "bar",
				Qos:       libmqtt.Qos1,
				Payload:   []byte("bar data"),
			},
			// ...
		)
	})

	client.UnSubscribe("foo", "bar")

	// wait until all connection exit
	client.Wait()
}
