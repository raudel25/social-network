port ?= 5000

.PHONY: dev
dev:
	go run cmd/main.go -port $(port)