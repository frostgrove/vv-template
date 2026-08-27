.PHONY: dev goose

CFG_PATH := config/app.yml

dev:
	go run cmd/main.go --config-path=$(CFG_PATH)

goose:
	go run cmd/goose/main.go --config-path=$(CFG_PATH) $(ARGS)
