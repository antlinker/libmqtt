# libmqtt command line app

A command line MQTT client built on top of libmqtt.

(please use it just for test purpose)

## Build

### Prerequisite

1. Go (with `GOPATH` configured)

### Steps

1. Go get and build this command line app

```bash
go get github.com/goiiot/libmqtt
cd $GOPATH/src/github.com/goiiot/libmqtt/cmd
go build -o libmqttc # or use `make` if you have make installed
```

2. Run client, and explore usages

```bash
./libmqttc # then type `h` or `help` for usage reference
```
