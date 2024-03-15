# Go Crud Example

- [Changelog](CHANGELOG.md)

---

## resources

- [go.dev](https://go.dev/) Go Programming Language
- [pgx](https://github.com/jackc/pgx) Replaces the standard `database/sql` driver
- [sqlx](https://github.com/jmoiron/sqlx) Extends the standard `database/sql` functions
- [echo](https://github.com/labstack/echo) minimalist web framework
- [goi18n](https://github.com/nicksnyder/go-i18n) internationalization
- [squirrel](https://github.com/Masterminds/squirrel) Query Builder
- [jsonSchema](https://github.com/santhosh-tekuri/jsonschema) Json Schema validation
- [echo-swagger](https://github.com/swaggo/echo-swagger) Swagger for echo

---

## Makefile

this project utilises docker to run `builds` and `migration`

`make [command]`

### Commands

- `init` install go development dependencies
- `build` build binary
- `run` run built binary
- `test` run tests
- `dev` build a docker image on local machine
- `dev/down` stops and remove dev docker image
- `migrate/up n=1` migrate database `n` steps
- `migrate/up.all` migrate database to latest
- `migrate/down n=1` rolls back `n` migrations
- `migrate/down.all` rolls back all migration
- `migration n=create_somethings_table` creates up and down sql in migrations
- `migrate/force n=23` force back failed migration version
- `refresh` runs down.all + up + seed
- `prune` prunes unused volumes, images and build caches
- `docker/ps` better format for docker ps command
- `audit` runs audit with go utilities on the project
- `list` lists update-able dependencies
- `update` downloads and updates project dependencies
- translations:
    - `translate.extract` update the `active.en.toml` file
    - `translate.merge` creates the `translate.ar.toml` file with new variables
    - translate the content of `translate.ar.toml` values
    - `translate.merge.done` merges translations to the `active.ar.toml` file
