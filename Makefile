.PHONY: dev generate goose

CFG_PATH := config/app.yml

dev:
	go run cmd/main.go --config-path=$(CFG_PATH)

generate:
	go run github.com/frostgrove/vv/cmd/vv generate -dir ./src/app

goose:
	go run cmd/goose/main.go --config-path=$(CFG_PATH) $(ARGS)
