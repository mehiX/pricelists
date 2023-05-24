.PHONY: test-unit
test-unit:
	go clean -testcache
	go test -race -cover -coverprofile=cover.out ./...
	go tool cover -html=cover.out -o cover.html

.PHONY: binary
binary:
	CGO_ENABLED=0 GOGC=off go build -o dist/prices ./cmd/prices/...