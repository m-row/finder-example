include .env
export

# Migrations
.PHONY: migrate.up migrate.up.all migrate.down migrate.down.all migrate.force
migrate/up:
	docker run --rm -v $(MIGRATIONS_ROOT):/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ -database $(CONNECTION_STRING) up $(n)
migrate/up/all:
	docker run --rm -v $(MIGRATIONS_ROOT):/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ -database $(CONNECTION_STRING) up
migrate/down:
	docker run --rm -v $(MIGRATIONS_ROOT):/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ -database $(CONNECTION_STRING) down $(n)
migrate/down/all:
	docker run --rm -v $(MIGRATIONS_ROOT):/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ -database $(CONNECTION_STRING) down -all
migration:
	docker run --rm -v $(MIGRATIONS_ROOT):/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ create -seq -ext=.sql -dir=./migrations $(n)
migrate/force:
	docker run --rm -v $(MIGRATIONS_ROOT):/migrations --network $(NETWORK) migrate/migrate -path=/migrations/ -database=$(CONNECTION_STRING) force $(n)

# Refresh
.PHONY: refresh
refresh: migrate/down/all migrate/up


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
