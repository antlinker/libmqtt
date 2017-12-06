
test:
	go test -v -run=. -bench=. -benchmem

lib:
	$(MAKE) -C c lib