BUILD_DIR = ../build
CFLAGS = -fPIC -O2 -std=c99 -I ../include -I $(BUILD_DIR) -L $(BUILD_DIR)
TARGET = $(BUILD_DIR)/example

LINUX_FLAGS = -lpthread

# detect build system for native build
HOST_SYSTEM = $(shell uname | cut -f 1 -d_)
SYSTEM ?= $(HOST_SYSTEM)

# linux
ifeq ($(SYSTEM),Linux)
	PLATFORM_FLAGS := -lpthread
endif
# macOS
ifeq ($(SYSTEM),Darwin)
	PLATFORM_FLAGS := -framework CoreFoundation -framework Security
endif

build:
	$(CC) $(CFLAGS) -o $(TARGET) example.c $(BUILD_DIR)/libmqtt.a $(PLATFORM_FLAGS)

run:
	$(TARGET)

clean:
	rm -rf $(TARGET)
