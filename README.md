# libmqtt

[![Build Status](https://travis-ci.org/goiiot/libmqtt.svg)](https://travis-ci.org/goiiot/libmqtt) [![GoDoc](https://godoc.org/github.com/goiiot/libmqtt?status.svg)](https://godoc.org/github.com/goiiot/libmqtt) [![GoReportCard](https://goreportcard.com/badge/goiiot/libmqtt)](https://goreportcard.com/report/github.com/goiiot/libmqtt)

MQTT 3.1.1 client lib with pure Go

## Contents

- [Features](#features)
- [Usage](#usage)
- [Topic Routing](#topic-routing)

## Features

1. A full functional MQTT 3.1.1 client. (currently only with in memory session state storage)
1. HTTP like API.
1. Customizable TopicRouter. (see [Topic Routing](#topic-routing))
1. Command line app support (see [cmd](./cmd/))
1. C lib support (see [c](./c/))
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

```go
client, err := libmqtt.NewClient(
    // server address(es)
    libmqtt.WithServer("localhost:1883"),
)
if err != nil {
    // handle client creation error
}
```

Notice: If you would like to explore all the options available, please refer to [godoc#Option](https://godoc.org/github.com/goiiot/libmqtt#Option)

4. Register the handlers and Connect, then you are ready to pub/sub with server

We recommend you to register handlers for pub, sub, unsub, net error and persist error, for they can provide you more controllability of the lifecycle of a MQTT client

```go
// register handler for pub success/fail (optional, but recommended)
client.HandlePub(PubHandler)

// register handler for sub success/fail (optional, but recommended)
client.HandleSub(SubHandler)

// register handler for unsub success/fail (optional, but recommended)
client.HandleUnSub(UnSubHandler)

// register handler for net error (optional, but recommended)
client.HandleNet(NetHandler)

// register handler for persist error (optional, but recommended)
client.HandlePersist(PersistHandler)

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
    // (the `server` is one of you have provided `servers` when create the client)
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

5. Unsubscribe topic(s).

```go
client.UnSubscribe("foo", "bar")
```

6. Destroy the client when you would like to.

```go
// passing true to Destroy means a immediate disconnect to server
// while passing false will try to send a DisConn packet to server
client.Destroy(true)
```

### As a C lib

Please refer to [c - README.md](./c/README.md).

### As a command line client

Please refer to [cmd - README.md](./cmd/README.md).

## Topic Routing

Routing topics is one of the most important thing when it comes to business logics, we currently have built two `TopicRouter`s which is ready to use, they are `TextRouter` and `RegexRouter`.

- `TextRouter` will match the exact same topic which was registered to client by `Handle` method. (this is the default router in a client).
- `RegexRouter` will go through all the registered topic handlers, and use regular expression to test whether that is matched and should dispatch to the handler.

If you would like to apply other routing strategy to the client, you can provide this option when creating the client

```go
client, err := NewClient(
    // for example, use `RegexRouter`
    libmqtt.WithRouter(NewRegexRouter()),
)
```

## RoadMap

1. File persist storage of session status. (High priority)
1. Full tested multiple connections in one client. (High priority)
1. More efficient processing. (Medium priority)
1. Add compatibility with mqtt 5.0 . (Medium priority)
1. Export to Java (JNI), Python (CPython), Objective-C... (Low priority)
