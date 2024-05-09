port ?= 8080

.PHONY: dev
dev:
	go run cmd/main.go -port $(port)