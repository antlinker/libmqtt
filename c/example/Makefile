BUILD_DIR = ../build
CFLAGS = -fPIC -O2 -std=c99 -I $(BUILD_DIR) -L $(BUILD_DIR)
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

# set ld path for linux, dyld path for darwin
# linux
ifeq ($(SYSTEM),Linux)
	RUNTIME_FLAGS := LD_LIBRARY_PATH=$(BUILD_DIR)
endif
# macOS
ifeq ($(SYSTEM),Darwin)
	RUNTIME_FLAGS := DYLD_LIBRARY_PATH=$(BUILD_DIR)
endif

.PHONY: build run clean

ARCHIVE_TARGET := $(TARGET)-a
SHARED_TARGET := $(TARGET)-so

build-single:
	$(CC) $(CFLAGS) -o $(ARCHIVE_TARGET) \
		example.c $(BUILD_DIR)/libmqtt.a $(PLATFORM_FLAGS)
	$(CC) $(CFLAGS) -o $(SHARED_TARGET) \
		-lmqtt $(PLATFORM_FLAGS) example.c

run: run-a

run-a:
	$(ARCHIVE_TARGET)

run-so:
	$(RUNTIME_FLAGS) $(SHARED_TARGET)

clean:
	rm -rf $(ARCHIVE_TARGET) $(SHARED_TARGET)
