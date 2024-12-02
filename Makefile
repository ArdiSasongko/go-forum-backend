build:
	@go build -o cmd/bin/main cmd/main.go

run:
	@./cmd/bin/main

tidy:
	@go mod tidy

.PHONY: build run tidy