# libmqtt command line client

A command line MQTT client built on top of libmqtt.

(please use it just for quick test purpose)

## Build

### Prerequisite

1. Go 1.9+ (with `GOPATH` configured)

### Steps

1. Go get and build this command line client

```bash
go get github.com/goiiot/libmqtt
cd $GOPATH/src/github.com/goiiot/libmqtt/cmd
go build -o libmqttc # or use `make` if you have make installed
```

2. Run, and explore usages

```bash
./libmqttc # then type `h` or `help` for usage reference
```
