BUILD_DIR := build
TARGET := libmqtt

.PHONY: lib clean example lib-shared lib-archive

lib: lib-shared lib-archive

lib-shared:
	env CGO_ENABLED=1 go build -buildmode=c-shared -ldflags=-s -o $(BUILD_DIR)/$(TARGET).so

lib-archive:
	env CGO_ENABLED=1 go build -buildmode=c-archive -ldflags=-s -o $(BUILD_DIR)/$(TARGET).a

example: lib
	$(MAKE) -C example build

clean:
	$(MAKE) -C example clean
	rm -rf $(BUILD_DIR)