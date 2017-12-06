BUILD_DIR := build
TARGET := libmqtt

.PHONY: lib clean test lib-shared lib-archive

lib: lib-shared lib-archive

lib-shared:
	env CGO_ENABLED=1 go build -buildmode=c-shared -o $(BUILD_DIR)/$(TARGET).so

lib-archive:
	env CGO_ENABLED=1 go build -buildmode=c-archive -o $(BUILD_DIR)/$(TARGET).a

test: lib
	$(MAKE) -c test test

clean:
	$(MAKE) -c test clean
	rm -rf $(build)