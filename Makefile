words-service:
	@go build -o bin/words-service cmd/main.go
	@./bin/words-service

docker-build:
	docker build -t word-service .

docker-run:
	docker run -p 3001:3001 word-service

test:
	go clean -testcache && go test ./internal/... && go test -coverprofile=coverage.out ./internal/...

.PHONY: words-service docker-build docker-run test-and-coverage