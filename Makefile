.PHONY: test lib clean

test:
	go test -v -run=. -bench=. -benchmem

lib:
	$(MAKE) -C c lib

clean:
	$(MAKE) -C c clean