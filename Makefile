# All tests must be run sequentially because the RTC can only be opened by one process.
# "-p 1" specifies parallelism, Ie. One test process.
.PHONY: test
test:
	go test -race -p 1 ./...
