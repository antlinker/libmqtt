# libmqtt Extension

Extensions for libmqtt client

## Extension List

- Persist Extension
    1. RedisPersist (Test) - Use redis as session state persist storage
- Router Extension
    1. HttpRouter (TODO) - HTTP path router for MQTT message

## Usage

1. Go get extension package

```bash
go get github.com/goiiot/libmqtt/extension
```

2. Import extensions

```go
import "github.com/goiiot/libmqtt/extension"
```