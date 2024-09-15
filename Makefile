include .env
export

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# Migrations
.PHONY: migrate.up migrate.up.all migrate.down migrate.down.all migrate.force
## migrate/up n=1: applies migration with argument (n) as steps count
migrate/up:
	docker run --rm -v $(MIGRATIONS_ROOT):/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ -database $(CONNECTION_STRING) up $(n)
## migrate/up/all: applies migrations up to latest
migrate/up/all:
	docker run --rm -v $(MIGRATIONS_ROOT):/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ -database $(CONNECTION_STRING) up
## migrate/down n=1: reverts migration with argument (n) as steps count
migrate/down:
	docker run --rm -v $(MIGRATIONS_ROOT):/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ -database $(CONNECTION_STRING) down $(n)
## migrate/down/all: reverts every migration down
migrate/down/all:
	docker run --rm -v $(MIGRATIONS_ROOT):/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ -database $(CONNECTION_STRING) down -all
## migration n=create_named_table: creates new up/down migration files with (n) as file name
migration:
	docker run --rm -v $(MIGRATIONS_ROOT):/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ create -seq -ext=.sql -dir=./migrations $(n)
## migrate/force n=1: fixes migration level if an error occured with (n) as the version number
migrate/force:
	docker run --rm -v $(MIGRATIONS_ROOT):/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ -database=$(CONNECTION_STRING) force $(n)

# Refresh
.PHONY: refresh
## refresh: applies migrate/down/all followed by migrate/up/all
refresh: migrate/down/all migrate/up


# Docker
.PHONY: prune docker/ps docker/inspect
## prune: clears unused docker data
prune:
	docker system prune -a -f --volumes
## docker/ps: formatted docker view of docker ps
docker/ps:
	docker ps --format "table {{.Names}}\t{{.Status}}\t{{.RunningFor}}\t{{.Size}}\t{{.Ports}}"
## docker/inspect n=name: inspect container (n) 
docker/inspect: 
	docker inspect -f "{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}" $(n)

# Go
.PHONY: list update
## list: golang packages
list:
	go list -m -u
## update: golang packages
update:
	go get -u ./...
