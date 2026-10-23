.PHONY: lint fmt test check

lint:
	golangci-lint run ./...

fmt:
	golangci-lint fmt

test:
	go test ./...

check: fmt lint test