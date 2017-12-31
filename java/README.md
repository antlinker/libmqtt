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

1. Create a client with builder

```java
LibMQTT.Client client = LibMQTT.newBuilder("localhost:8883")
    .setCleanSession(true)
    .setDialTimeout(10)
    .setKeepalive(10, 1.2)
    .setIdentity(sUsername, sPassword)
    .setLog(LibMQTT.LogLevel.Verbose)
    .setSendBuf(100).setRecvBuf(100)
    .setTLS(sClientCert, sClientKey, sCACert, sServerName, true)
    .setClientID(sClientID)
    .build();
```

2. Set the client lifecycle callback

```java
client.setCallback(new LibMQTT.Callback() {
    public void onConnected() {
        println("connected to server");
            client.subscribe(sTopicName, 0);
    }

    public void onSubResult(String topic, String err){
        if (err != null) {
            println("sub", topic ,"failed:", err);
        }
        println("sub", topic ,"success");
        client.publish(sTopicName, 0, sTopicMsg.getBytes());
    }

    public void onPubResult(String topic, String err){
        if (err != null) {
            println("pub", topic ,"failed:", err);
        }
        println("pub", topic ,"success");
        client.unsubscribe(sTopicName);
    }

    public void onUnSubResult(String topic, String err){
        if (err != null) {
            println("unsub", topic ,"failed:", err);
        }
        println("unsub", topic ,"success");
        client.destroy(true);
    }

    public void onPersistError(String err) {
        println("persist err happened:", err);
    }

    public void onLost(String err) {
        println("connection lost, err:", err);
    }
});
```

3. Connect to server

For java doesn't allow jni part to lock the main thread, you have to implement your own method to wait for client exit, here we just make use of `Thread.sleep()` for example

```java
client.connect();
// Thread.sleep(100 *1000);
```