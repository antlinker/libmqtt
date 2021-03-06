# libmqtt

[![Build Status](https://travis-ci.org/goiiot/libmqtt.svg)](https://travis-ci.org/goiiot/libmqtt) [![GoDoc](https://godoc.org/github.com/goiiot/libmqtt?status.svg)](https://godoc.org/github.com/goiiot/libmqtt) [![GoReportCard](https://goreportcard.com/badge/goiiot/libmqtt)](https://goreportcard.com/report/github.com/goiiot/libmqtt)

Feature rich modern MQTT 3.1.1 client lib in pure Go, for `Go`, `C/C++`, `Java` and `Python`

## Contents

- [Features](#features)
- [Extensions](#extensions)
- [Usage](#usage)
- [Topic Routing](#topic-routing)
- [Session Persist](#session-persist)
- [Benchmark](#benchmark)
- [RoadMap](#roadmap)
- [LICENSE](#license)

## Features

1. Feature rich MQTT 3.1.1 client
1. HTTP server like API
1. High performance and less memory footprint (see [Benchmark](#benchmark))
1. Customizable `TopicRouter` (see [Topic Routing](#topic-routing))
1. Builtin multiple session persist methods (see [Session Persist](#session-persist))
1. [C/C++ lib](./c/), [Java lib](./java/), [Python lib - TODO](./python/), [Command line client](./cmd/) support
1. Idiomatic Go, reactive stream

## Extensions

Helpful extensions for libmqtt (see [extension](./extension/))

## Usage

This project can be used as

- A [Go lib](#as-a-go-lib)
- A [C/C++ lib](#as-a-c-lib)
- A [Java lib](#as-a-java-lib)
- A [Python lib](#as-a-python-lib) (TODO)
- A [Command line client](#as-a-command-line-client)

### As a Go lib

#### Prerequisite

- Go 1.9+ (with `GOPATH` configured)

#### Steps

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

Notice: If you would like to explore all the options available, please refer to [GoDoc#Option](https://godoc.org/github.com/goiiot/libmqtt#Option)

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

5. Unsubscribe topic(s)

```go
client.UnSubscribe("foo", "bar")
```

6. Destroy the client when you would like to

```go
// passing true to Destroy means a immediate disconnect to server
// while passing false will try to send a DisConn packet to server
client.Destroy(true)
```

### As a C/C++ lib

Please refer to [c - README.md](./c/README.md)

### As a Java lib

Please refer to [java - README.md](./java/README.md)

### As a Python lib

TODO

### As a command line client

Please refer to [cmd - README.md](./cmd/README.md)

## Topic Routing

Routing topics is one of the most important thing when it comes to business logics, we currently have built two `TopicRouter`s which is ready to use, they are `TextRouter` and `RegexRouter`

- `TextRouter` will match the exact same topic which was registered to client by `Handle` method. (this is the default router in a client)
- `RegexRouter` will go through all the registered topic handlers, and use regular expression to test whether that is matched and should dispatch to the handler

If you would like to apply other routing strategy to the client, you can provide this option when creating the client

```go
client, err := libmqtt.NewClient(
    // ...
    // for example, use `RegexRouter`
    libmqtt.WithRouter(libmqtt.NewRegexRouter()),
    // ...
)
```

## Session Persist

Per MQTT Specification, session state should be persisted and be recovered when next time connected to server without clean session flag set, currently we provide persist method as following:

1. `NonePersist` - no session persist
1. `MemPersist` - in memory session persist
1. `FilePersist` - files session persist (with write barrier)
1. `RedisPersist` - redis session persist (available inside [github.com/goiiot/libmqtt/extension](./extension/) package)

__Note__: Use `RedisPersist` if possible.

## Benchmark

The procedure of the benchmark is as following:

1. Create the client
1. Connect to server
1. Subscribe to topic `foo`
1. Publish to topic `foo`
1. Unsubsecibe when received all published message (with `foo` topic)
1. Destroy client (a sudden disconnect without disconnect packet)

The benchmark result listed below was taken on a Macbook Pro 13' (Early 2015, macOS 10.13.2), statistics inside which is the value of ten times average

|Bench Name|Pub Count|ns/op|B/op|allocs/op|Transfer Time|Total Time|
|---|---|---|---|---|---|---|
|BenchmarkPahoClient-4|10000|199632|1399|31|0.230s|2.021s|
|BenchmarkLibmqttClient-4|10000|144407|331|9|0.124s|1.467s|
|BenchmarkPahoClient-4|50000|205884|1395|31|1.170s|10.316s|
|BenchmarkLibmqttClient-4|50000|161640|328|9|0.717s|8.105s|

You can make the benchmark using source code from [benchmark](./benchmark/)

Notice: benchmark on libmqtt sometimes can be a infinite loop, we are now trying to solve that

## RoadMap

1. File persist storage of session status (High priority)
1. Full tested multiple connections in one client (High priority)
1. Add compatibility with mqtt 5.0 (Medium priority)
1. Export to Python (CPython)... (Low priority)

## LICENSE

```text
Copyright GoIIoT (https://github.com/goiiot)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```