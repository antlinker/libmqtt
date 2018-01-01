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