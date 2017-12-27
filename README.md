# libmqtt

[![Build Status](https://travis-ci.org/goiiot/libmqtt.svg)](https://travis-ci.org/goiiot/libmqtt) [![GoDoc](https://godoc.org/github.com/goiiot/libmqtt?status.svg)](https://godoc.org/github.com/goiiot/libmqtt) [![GoReportCard](https://goreportcard.com/badge/goiiot/libmqtt)](https://goreportcard.com/report/github.com/goiiot/libmqtt)

MQTT 3.1.1 client lib with pure Go

(and will be exported to C soon)

## Features

1. A full functional MQTT 3.1.1 client (currently without session state storage)
1. HTTP server like API
1. Customizable TopicRouter (predefined `TextRouter`, `RegexRouter` and `RestRouter` in work)
1. More efficient, idiomatic Go

## Usage

### As a Go lib

1. Go get this package

```go
go get github.com/goiiot/libmqtt
```

2. Import this package in your project file

```go
import "github.com/goiiot/libmqtt"
```

3. Create a custom client

If you would like to explore all the options avaliable, please refer to [godoc](https://godoc.org/github.com/goiiot/libmqtt)

```go
client := libmqtt.NewClient(
    // server address(es)
    libmqtt.WithServer("localhost:1883"),
)
```

4. Register the handlers and Connect, then you are ready to pub/sub with server

```go
// define your topic handlers like a golang http server
client.Handle("SomeTopic", func(topic string, qos libmqtt.QosLevel, msg []byte) {
    // handle the topic message
})

client.Handle("ExampleTopic", func(topic string, qos libmqtt.QosLevel, msg []byte) {
    // handle the topic message
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
        &libmqtt.Topic{Name: "SomeTopic"},
        &libmqtt.Topic{Name: "ExampleTopic", Qos: libmqtt.Qos1},
        // ...
    )

    // publish some topic message(s)
    client.Publish(
        &libmqtt.PublishPacket{TopicName: "SomeTopic", Qos: libmqtt.Qos0, Payload: []byte("something bad")},
        &libmqtt.PublishPacket{TopicName: "ExampleTopic", Qos: libmqtt.Qos1, Payload: []byte("something good")},
        // ...
    )
})
```

5. Unsubscribe topic(s)

```go
client.UnSubscribe("foo", "bar")
```

## TODO

1. persistant storage of session status
1. full tested multiple connections in one client
1. commandline app
1. export to C, Java, Objective-C, Python...
1. more efficient processing