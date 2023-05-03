.PHONY: build lint vet test check

name = map_builder

build: $(wildcard $(shell find . -type f | grep "\.go"))
	CGO_ENABLED=0 go build -o ./$(name) -a -installsuffix cgo -ldflags '-s' .

.PHONY: test
test:
	go test -cover `go list ./... | grep -v /vendor/`
