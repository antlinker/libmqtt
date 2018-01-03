# libmqtt C/C++ lib

A MQTT C/C++ client library exported from libmqtt

__NOTE__: This library is still under work, some bugs can happen

## Contents

1. [Build](#build)
1. [Usage](#usage)

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

0. Include the generated header `libmqtt.h`

```c
#include "libmqtt.h"
```

1. Create and configure the client

```c
client = Libmqtt_new_client();
if (client < 1) {
  printf("create client failed\n");
  exit(1);
}

Libmqtt_client_set_server(client, "localhost:8883");
Libmqtt_client_set_log(client, libmqtt_log_silent);
Libmqtt_client_set_clean_session(client, true);
Libmqtt_client_set_keepalive(client, 10, 1.2);
Libmqtt_client_set_identity(client, "foo", "bar");
Libmqtt_client_set_client_id(client, "foo");
Libmqtt_client_set_tls(client, "client-cert.pem", "client-key.pem", "ca-cert.pem", "MacBook-Air.local", true);

// other options provided by libmqtt are also avalieble here
```

2. Setup the client and check if error happened

```c
char *err = Libmqtt_setup(client);
if (err != NULL) {
  printf("error happened when create client: %s\n", err);
  return 1;
}
```

3. Create and register the handlers as you like

```c
// sub_handler for subscribe packet response, non-null err means topic sub failed
void sub_handler(char *topic, int qos, char *err) {}

// pub_handler for publish packet response, non-null err means `topic` publish failed
void pub_handler(char *topic, char *err) {}

// unsub_handler for unsubscribe response, non-null err means `topic` unsubscribe failed
void unsub_handler(char *topic, char *err) {}

// net_handler for net connection information, called only when error happens
void net_handler(char *server, char *err) {}

// topic_handler, just a example topic handler
void topic_handler(char *topic, int qos, char *msg, int size) {}

// client handlers (optional, but recommend to)
Libmqtt_set_pub_handler(client, &pub_handler);
Libmqtt_set_sub_handler(client, &sub_handler);
Libmqtt_set_net_handler(client, &net_handler);
Libmqtt_set_unsub_handler(client, &unsub_handler);

// the topic handler
Libmqtt_handle(client, "foo", &topic_handler);
// Libmqtt_handle(client, "bar", &topic_handler);
// ...
```

4. Connect to server and wait for connection close

```c
// conn_handler for connect packet response
void conn_handler(char *server, libmqtt_connack_t code, char *err) {}

// connect to server with
Libmqtt_conn(client, &conn_handler);

// wait until all connection closed
Libmqtt_wait(client);
```

You can refer to the [example](./example/) for a full usage example

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