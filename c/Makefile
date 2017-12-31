BUILD_DIR := build
TARGET := libmqtt
BUILD_ARGS := -ldflags=-s

LIB_SHARED := $(BUILD_DIR)/$(TARGET).so
LIB_ARCHIVE := $(BUILD_DIR)/$(TARGET).a

.PHONY: lib example clean

lib: $(LIB_ARCHIVE) $(LIB_SHARED)

$(LIB_SHARED):
	env CGO_ENABLED=1 go build -buildmode=c-shared \
	    $(BUILD_ARGS) -o $(LIB_SHARED)

$(LIB_ARCHIVE):
	env CGO_ENABLED=1 go build -buildmode=c-archive \
	    $(BUILD_ARGS) -o $(LIB_ARCHIVE)

example: $(LIB_ARCHIVE)
	$(MAKE) -C example build

clean:
	$(MAKE) -C example clean
	rm -rf $(BUILD_DIR)