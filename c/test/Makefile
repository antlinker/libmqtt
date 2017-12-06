
TARGETS = conn_test ping_test pub_test sub_test

CFLAGS = -Wall -fPIC -O2 -std=c99 \
		 -I../include -I../build

LDFLAGS = -L../build -lmqtt

%.o: %.c
	$(CC) $(CFLAGS) -c -o $@ $<

ping_test: ping_test.o
	$(CC) $(CFLAGS) $(LDFLAGS) -o $@ $^

clean:
	rm -rf *.o
	rm -rf $(TARGETS)
