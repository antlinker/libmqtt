# libmqtt

[![Build Status](https://travis-ci.org/goiiot/libmqtt.svg)](https://travis-ci.org/goiiot/libmqtt) [![GoDoc](https://godoc.org/github.com/goiiot/libmqtt?status.svg)](https://godoc.org/github.com/goiiot/libmqtt) [![GoReportCard](https://goreportcard.com/badge/goiiot/libmqtt)](https://goreportcard.com/report/github.com/goiiot/libmqtt)

MQTT 3.1.1 client lib with pure Go

## Features

1. A full functional MQTT 3.1.1 client (currently without session state storage)
1. HTTP server like API
1. Customizable TopicRouter (predefined `TextRouter`, `RegexRouter` and `RestRouter` in work)
1. Command line app support (see [cmd](./cmd/))
1. Exported to C lib (see [c](./c/))
1. More efficient, idiomatic Go (maybe not quite idiomatic for now)

## Usage

This project can be used as

- A [Go lib](#as-a-go-lib)
- A [C lib](#as-a-c-lib)
- A [Command line client](#as-a-command-line-client)

### As a Go lib

1. Go get this project

```bash
go get github.com/goiiot/libmqtt
```

2. Import this package in your project file

```go
import "github.com/goiiot/libmqtt"
```

3. Create a custom client

If you would like to explore all the options available, please refer to [godoc](https://godoc.org/github.com/goiiot/libmqtt)

```go
client := libmqtt.NewClient(
    // server address(es)
    libmqtt.WithServer("localhost:1883"),
)
```

4. Register the handlers and Connect, then you are ready to pub/sub with server

```go
// define your topic handlers like a golang http server
client.Handle("foo", func(topic string, qos libmqtt.QosLevel, msg []byte) {
    // handle the topic message
})

client.Handle("bar", func(topic string, qos libmqtt.QosLevel, msg []byte) {
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
        }, &libmqtt.PublishPacket{
            TopicName: "bar",
            Qos:       libmqtt.Qos1,
            Payload:   []byte("bar data"),
        },
        // ...
    )
})
```

5. Unsubscribe topic(s)

```go
client.UnSubscribe("foo", "bar")
```

### As a C lib

Please refer to [c - README.md](./c/README.md)

### As a command line client

Please refer to [cmd - README.md](./cmd/README.md)

## RoadMap

1. Persistant storage of session status (High priority)
1. Full tested multiple connections in one client (High priority)
1. More efficient processing (Medium priority)
1. Add compatibility with mqtt 5.0 (Medium priority)
1. Export to Java (JNI), Python (CPython), Objective-C... (Low priority)
