# libmqtt Java lib

A MQTT Java client library exported from [libmqtt](https://github.com/goiiot/libmqtt) with JNI

__NOTE__: This library is still under work, some bugs can happen

## Contents

1. [Build](#build)
1. [Usage](#usage)

## Build

### Prerequisite

1. Go (with `GOPATH` configured)
1. JDK 1.6+
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

0. Import package

```java
import cc.goiiot.libmqtt.*;
```

1. Create a client with client builder

```java
Client client = Client.newBuilder("localhost:8883")
    .setTLS("client-cert.pem", "client-key.pem", "ca-cert.pem", "server.name", true)
    .setCleanSession(true)
    .setDialTimeout(10)
    .setKeepalive(10, 1.2)
    .setIdentity("foo", "bar")
    .setLog(LogLevel.Verbose)
    .setSendBuf(100)
    .setRecvBuf(100)
    .setClientID("foo")
    .build();
```

2. Set the client lifecycle callback

```java
client.setCallback(new Callback() {
    public void onConnResult(boolean ok, String descp) {
        if (!ok) {
            println("connection error:", descp);
            return;
        }
        println("connected to server");
        client.subscribe(sTopicName, 0);
    }

    public void onLost(String descp) {
        println("connection lost, err:", descp);
    }

    public void onSubResult(String topic, boolean ok, String descp) {
        if (!ok) {
            println("sub", topic, "failed:", descp);
            return;
        }
        println("sub", topic, "success");
        client.publish(sTopicName, 0, sTopicMsg.getBytes());
    }

    public void onPubResult(String topic, boolean ok, String descp) {
        if (!ok) {
            println("pub", topic, "failed:", descp);
            return;
        }
        println("pub", topic, "success");
        client.unsubscribe(sTopicName);
    }

    public void onUnSubResult(String topic, boolean ok, String descp) {
        if (!ok) {
            println("unsub", topic, "failed:", descp);
            return;
        }
        println("unsub", topic, "success");
        client.destroy(true);
    }

    public void onPersistError(String descp) {
        println("persist err happened:", descp);
    }
});
```

3. Handle the topic message

__NOTE__: Current routing behavior of the Java lib is to route all topic message to the last registered topic handler, if that's null, no message will be notified, we are now working on it)

```java
client.handle(sTopicName, new TopicMessageCallback(){
    public void onMessage(String topic, int qos, byte[] payload) {
        println("received message:", new String(payload), "from topic:", topic, "qos = " + qos);
    }
});
```

4. Connect to server and wait for exit

__NOTE__: Currenly, once called `waitClient()` will not exit even when all connection has been closed, we are working on that.

```java
client.connect();
client.waintClient();
```

You can refer to the [example](./cc/goiiot/libmqtt/example/) for a full usage example

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