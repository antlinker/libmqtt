# libmqtt C lib

A MQTT C client library exported from libmqtt

## Temporary Restrictions

1. Only one client, and one connection per client supported

## Build

### Prerequisite

1. Go (with `GOPATH` configured)
1. make (for build ease)
1. gcc (for cgo builds)

### Steps

1. Go get this project

```bash
go get github.com/goiiot/libmqtt
```

2. Make the library

```bash
cd $GOPATH/src/github.com/goiiot/libmqtt/c

make
```

## Usage

1. Configure the client

```c
SetServer(SERVER);              // required, or client will not work
SetLog(libmqtt_log_verbose);    // open log as you wish
SetCleanSession(true);          // clean seesion mark
SetKeepalive(10, 1.2);          // set keepalive options
```

2. Setup the client

```c
SetUp();
```

3. Register the handlers as you like

```c
// client handlers (optional, but recommend to)
SetConnHandler(&conn_handler);
SetPubHandler(&pub_handler);
SetSubHandler(&sub_handler);
SetUnSubHandler(&unsub_handler);
SetNetHandler(&net_handler);

// the topic handler
Handle("foo", &topic_handler);
// Handle("bar", &topic_handler);
// ...
```

4. Connect to server and wait for close

```c
// connect to server with
Conn();

// wait until all connection closed
Wait();
```

You can refer to the [example](./example/) for a full usage example
