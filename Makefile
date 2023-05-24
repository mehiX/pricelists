.PHONY: test
test:
	docker compose up -d --build database
	go clean -testcache
	go test -v -race -cover -coverprofile=cover.out ./...
	go tool cover -html=cover.out -o cover.html
	docker compose down --volumes --remove-orphans

.PHONY: binary
binary:
	CGO_ENABLED=0 GOGC=off go build -o dist/prices ./cmd/prices/...

.PHONY: up
up:
	docker compose up -d --build

.PHONY: db
db:
	docker compose up -d --build database

.PHONY: logs
logs:
	docker compose logs -f

.PHONY: down
down:
	docker compose down --volumes --remove-orphans