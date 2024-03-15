include .env
export

.PHONY: init
 init:
	@go install github.com/cosmtrek/air@latest
	@go install github.com/go-delve/delve/cmd/dlv@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.1
	@go install github.com/jesseduffield/lazygit@latest
	@go install github.com/mfridman/tparse@latest
	@go install github.com/nametake/golangci-lint-langserver@latest
	@go install github.com/nicksnyder/go-i18n/v2/goi18n@latest
	@go install github.com/segmentio/golines@latest
	@go install honnef.co/go/tools/cmd/staticcheck@latest
	@go install mvdan.cc/gofumpt@latest

# Application
#
# Translations
.PHONY: i18n/extract i18n/merge
i18n/extract:
	@goi18n extract --outdir .
i18n/merge: 
	@goi18n merge --outdir . ./active.ar.toml ./active.en.toml
i18n/done:
	@goi18n merge --sourceLanguage ar --outdir . ./active.ar.toml ./translate.ar.toml

# Migrations
.PHONY: migrate.up migrate.up.all migrate.down migrate.down.all migrate.force
migrate.up:
	docker run --rm -v $(MIGRATIONS_ROOT):/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ -database $(CONNECTION_STRING) up $(n)
migrate.up.all:
	docker run --rm -v $(MIGRATIONS_ROOT):/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ -database $(CONNECTION_STRING) up
migrate.down:
	docker run --rm -v $(MIGRATIONS_ROOT):/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ -database $(CONNECTION_STRING) down $(n)
migrate.down.all:
	docker run --rm -v $(MIGRATIONS_ROOT):/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ -database $(CONNECTION_STRING) down -all
migration:
	docker run --rm -v $(MIGRATIONS_ROOT):/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ create -seq -ext=.sql -dir=./migrations $(n)
migrate.force:
	docker run --rm -v $(MIGRATIONS_ROOT):/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ -database=$(CONNECTION_STRING) force $(n)

# Refresh
.PHONY: refresh
refresh: migrate.down.all migrate.up


# Docker
.PHONY: prune docker.ps docker.inspect
prune:
	docker system prune -a -f --volumes
docker.ps:
	docker ps --format "table {{.Names}}\t{{.Status}}\t{{.RunningFor}}\t{{.Size}}\t{{.Ports}}"
# inspect a container local ip n=name of container
docker.inspect: 
	docker inspect -f "{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}" $(n)

# Go
.PHONY: list update
list:
	go list -m -u
update:
	go get -u ./...

# ==================================================================================== #
# QUALITY CONTROL
# ==================================================================================== #
## audit: tidy dependencies and format, vet and test all code
.PHONY: audit
audit:
	@echo 'Tidying and verifying module dependencies...'
	go mod tidy
	go mod verify
	@echo 'Formatting code...'
	gofumpt -l -w -extra .
	@echo 'Vetting code...'
	go vet ./...
	staticcheck ./...
	@echo 'Running tests...'
	go test -race -vet=off ./...

