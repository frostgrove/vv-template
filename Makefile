CFG_PATH := config/app.yml
# Everything after `make migrate` is handed to the vvgoose CLI: `make migrate status`,
# `make migrate table users`, `make migrate rollback 2`. Bare `make migrate` applies
# pending migrations; open the interactive menu with `go run ./cmd/migrate`.
ifneq ($(filter migrate goose,$(MAKECMDGOALS)),)
MIGRATE_ARGS := $(filter-out migrate goose,$(MAKECMDGOALS))
endif

.PHONY: dev seed migrate goose generate generate-check build vet lint test up down workspace

dev:
	go run ./cmd --config-path=$(CFG_PATH)

seed:
	go run ./cmd/seed --config-path=$(CFG_PATH)

migrate:
	go run ./cmd/migrate --config-path=$(CFG_PATH) $(if $(MIGRATE_ARGS),$(MIGRATE_ARGS),migrate)

# The old name of `migrate`.
goose: migrate

generate:
	go run github.com/frostgrove/vv/cmd/vv generate -dir ./src/mod

generate-check:
	go run github.com/frostgrove/vv/cmd/vv generate -dir ./src/mod -check

build:
	go build ./...

vet:
	go vet ./...

lint:
	golangci-lint run ./...

test:
	go test ./...

# Build against a local framework checkout instead of the pinned commit.
FRAMEWORK ?= ../framework
FRAMEWORK_MODULES := . app/appfx app/http/appfiber auth/http/authfiber crud/adapter/crudsql/crudsqlfx \
	crud/http/crudfiber health/healthfx otel runtime/runtimefx utils/vvcfg utils/vvgoose

workspace:
	rm -f go.work go.work.sum
	go work init . $(addprefix $(FRAMEWORK)/,$(FRAMEWORK_MODULES))

up:
	docker compose -f compose.dev.yml up -d --wait

down:
	docker compose -f compose.dev.yml down

ifneq ($(MIGRATE_ARGS),)
$(MIGRATE_ARGS):
	@:
endif
