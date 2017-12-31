# libmqtt Java lib

A MQTT Java client library exported from [libmqtt](https://github.com/goiiot/libmqtt) with JNI

__NOTE__: This library is still under work, some bugs can happen

## Contents

1. [Build](#build)
1. [Usage](#usage)

## Build

### Prerequisite

1. Go (with `GOPATH` configured)
1. JDK 1.7+
1. make (for build ease)
1. gcc (for cgo builds)

### Steps

1. Go get this project

```bash
go get github.com/goiiot/libmqtt
```

2. Make the basic C library

```bash
cd $GOPATH/src/github.com/goiiot/libmqtt/c

make
```

3. Make the Java library

```bash
cd ../java

make JAVA_HOME=/path/to/your/java/home
```

4. Run the Example and explore how to make it work

```bash
make run-jni
```


## Usage

0. Include the predefined header `libmqtth.h` and generated header `libmqtt.h`

```c
#include "libmqtth.h"
#include "libmqtt.h"
```

1. Configure the client

```c
SetServer("localhost:1883");    // required, or client will not work
SetLog(libmqtt_log_verbose);    // open log as you wish
SetCleanSession(true);          // clean seesion mark
SetKeepalive(10, 1.2);          // set keepalive options
// other options provided by libmqtt are also avalieble here
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

4. Connect to server and wait for connection close

```c
// connect to server with
Conn();

// wait until all connection closed
Wait();
```

You can refer to the [example](./example/) for a full usage example
