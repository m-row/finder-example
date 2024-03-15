package models

import (
	"github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/m-row/finder-example/model"
	"github.com/m-row/finder-example/models/permission"
	"github.com/m-row/finder-example/models/role"
	"github.com/m-row/finder-example/models/user"
)

type Models struct {
	DB *sqlx.DB
	QB *squirrel.StatementBuilderType

	User       *user.Queries
	Role       *role.Queries
	Permission *permission.Queries
}

func Setup(
	db *sqlx.DB,
	info map[string][]string,
) *Models {
	dbCache := squirrel.NewStmtCache(db)

	qb := squirrel.
		StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).
		RunWith(dbCache)

	d := &model.Dependencies{
		DB:     db,
		QB:     &qb,
		PGInfo: info,
	}

	return &Models{
		DB: db,
		QB: &qb,

		User:       user.New(d),
		Role:       role.New(d),
		Permission: permission.New(d),
	}
}
