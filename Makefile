GOOSE_VERSION := v3.11.2
POSTGRES := postgresql://postgres:postgres@postgres:5432/postgres
MIGRATIONS_DIR := deployments/migrations

.PHONY: generate
generate:
	go get golang.org/x/tools/go/packages
	go get golang.org/x/tools/go/ast/astutil
	go get golang.org/x/tools/imports
	go get github.com/urfave/cli/v2
	go run github.com/99designs/gqlgen generate


.PHONY: install-goose
install-goose:
	@go install github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION)

.PHONY: migrate-run
migrate-run:
	@goose --dir=$(MIGRATIONS_DIR) postgres "$(POSTGRES)" up

migrate-down:
	@goose --dir=$(MIGRATIONS_DIR) postgres "$(POSTGRES)" down

MIGRATION := table_init
.PHONY: new-migration
new-migration:
	goose --dir="$(MIGRATIONS_DIR)" create $(MIGRATION) sql


CONFIG_FILE=./config/config.yaml

in-memory-docker:
	sed -i 's/DB: ".*"/DB: "in_memory"/' $(CONFIG_FILE)
	docker compose -f docker-compose.yml --profile in-memory up

postgres-docker:
	sed -i 's/DB: ".*"/DB: "postgres"/' $(CONFIG_FILE)
	sed -i 's/host:.* ".*"/host: "postgres"/' $(CONFIG_FILE)
	docker compose -f docker-compose.yml --profile postgres up

in-memory-local:
	sed -i 's/DB: ".*"/DB: "in_memory"/' $(CONFIG_FILE)
	go run ./cmd/app/main.go


postgres-local:
	sed -i 's/DB: ".*"/DB: "postgres"/' $(CONFIG_FILE)
	sed -i 's/host:.* ".*"/host: "localhost"/' $(CONFIG_FILE)
	go run ./cmd/app/main.go
