.PHONY: test lib client clean fuzz-test

test:
	go test -v -run=. -bench=. -benchmem

lib:
	$(MAKE) -C c lib

client:
	$(MAKE) -C cmd build

clean: fuzz-clean
	$(MAKE) -C c clean
	$(MAKE) -C cmd clean

fuzz-test:
	go-fuzz-build github.com/goiiot/libmqtt
	go-fuzz -bin=./libmqtt-fuzz.zip -workdir=fuzz-test

fuzz-clean:
	rm -rf fuzz-test libmqtt-fuzz.zip
