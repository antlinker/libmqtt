C_BUILD_DIR := ../c/build
BUILD_DIR := build

JAVA_PACKAGE := cc.goiiot.libmqtt
JAVA_PACKAGE_PATH := cc/goiiot/libmqtt

JAVA_HOME ?= /Library/Java/JavaVirtualMachines/jdk1.8.0_131.jdk/Contents/Home
CFLAGS := -I $(C_BUILD_DIR) -I build -L $(C_BUILD_DIR)

SRC_SET := libmqtt_jni.c

# detect build system for native build
HOST_SYSTEM = $(shell uname | cut -f 1 -d_)
SYSTEM ?= $(HOST_SYSTEM)

# linux
ifeq ($(SYSTEM),Linux)
	PLATFORM_FLAGS := -fPIC -shared -lpthread -I$(JAVA_HOME)/include -I$(JAVA_HOME)/include/linux
	SUFFIX=so
endif
# macOS
ifeq ($(SYSTEM),Darwin)
	PLATFORM_FLAGS := -dynamiclib -framework JavaVM -framework CoreFoundation -framework Security -I$(JAVA_HOME)/include -I$(JAVA_HOME)/include/darwin
	SUFFIX=jnilib
endif

.PHONY: build compile-java gen-jni-header clean

build: gen-jni-header
	$(CC) $(CFLAGS) $(PLATFORM_FLAGS) -o $(BUILD_DIR)/libmqtt-jni.$(SUFFIX) -lmqtt $(SRC_SET)

gen-jni-header: compile-java
	javah -jni -classpath $(BUILD_DIR) -d $(BUILD_DIR) $(JAVA_PACKAGE).LibMQTT

compile-java:
	mkdir -p build
	javac $(JAVA_PACKAGE_PATH)/LibMQTT.java -d $(BUILD_DIR)

run-example: compile-example
	cd $(BUILD_DIR) && \
	DYLD_LIBRARY_PATH=../$(C_BUILD_DIR) \
	LD_LIBRARY_PATH=../$(C_BUILD_DIR) \
	java -Djava.library.path=. $(JAVA_PACKAGE_PATH)/example/Example && \
	cd ..

compile-example:
	javac -Xdiags:verbose $(JAVA_PACKAGE_PATH)/example/Example.java -d $(BUILD_DIR)

clean:
	rm -rf $(BUILD_DIR)