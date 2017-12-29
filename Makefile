.PHONY: test lib client clean

test:
	go test -v -run=. -bench=. -benchmem

lib:
	$(MAKE) -C c lib

client:
	$(MAKE) -C cmd build

clean:
	$(MAKE) -C c clean
	$(MAKE) -C cmd clean