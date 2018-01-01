.PHONY: test lib client clean fuzz-test

test:
	go test -v -run=.

.PHONY: all-lib c-lib java-lib py-lib \
		clean-c-lib clean-java-lib clean-py-lib

all-lib: all-lib c-lib java-lib py-lib

clean-all-lib: clean-c-lib clean-java-lib clean-py-lib

c-lib:
	$(MAKE) -C c lib

clean-c-lib:
	$(MAKE) -C c clean

java-lib:
	$(MAKE) -C java lib

clean-java-lib:
	$(MAKE) -C java clean

py-lib:
	$(MAKE) -C python lib

clean-py-lib:
	$(MAKE) -C python clean

client:
	$(MAKE) -C cmd build

clean: clean-all-lib fuzz-clean

fuzz-test:
	go-fuzz-build github.com/goiiot/libmqtt
	go-fuzz -bin=./libmqtt-fuzz.zip -workdir=fuzz-test

fuzz-clean:
	rm -rf fuzz-test libmqtt-fuzz.zip
