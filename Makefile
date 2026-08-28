CFG_PATH := config/app.yml
GOOSE_ARGS := $(filter-out dev generate goose,$(MAKECMDGOALS))

.PHONY: dev generate goose $(GOOSE_ARGS)

dev:
	go run cmd/main.go --config-path=$(CFG_PATH)

generate:
	go run github.com/frostgrove/vv/cmd/vv generate -dir ./src/app

goose:
	go run cmd/goose/main.go --config-path=$(CFG_PATH) $(GOOSE_ARGS)

$(GOOSE_ARGS):
	@: