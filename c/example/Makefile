BUILD_DIR = ../build
CFLAGS = -fPIC -O2 -std=c99 -I ../include -I $(BUILD_DIR)
TARGET = $(BUILD_DIR)

example:
	$(CC) $(CFLAGS) -o $(TARGET) example.c $(BUILD_DIR)/libmqtt.a

clean:
	rm -rf $(TARGET)
