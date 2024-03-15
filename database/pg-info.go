package database

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type tableColumn struct {
	Table      string `db:"table_name"`
	ColumnName string `db:"column_name"`
}

// PGInfo preloads all database meta table and column names sorted, into a
// map[string][]string where key is table_name and values are column names
// used in meta.columns object
//
//	banners             | city_id
//	banners             | created_at
//	banners             | id
//	banners             | img
//	banners             | is_disabled
//	banners             | sort
//	banners             | thumb
//	banners             | updated_at
//	banners             | url
//	banners             | vendor_id
//	cart_items          | cart_id
//	cart_items          | item_id
//	cart_items          | quantity
func PGInfo(conn *sqlx.DB, pgInfo map[string][]string) error {
	query := `
        SELECT
            table_name,
            column_name
        FROM
            information_schema.columns
        WHERE table_schema = 'public' AND
            table_name NOT IN (
                'geography_columns',
                'geometry_columns',
                'spatial_ref_sys',
                'schema_migrations'
            ) AND
            "column_name" NOT IN (
                'password_hash'
            )
        ORDER BY "table_name" ASC, "column_name" ASC
    `
	tc := []tableColumn{}
	if err := conn.SelectContext(
		context.Background(),
		&tc,
		query,
	); err != nil {
		return err
	}
	for i := range tc {
		r := tc[i]
		if _, found := pgInfo[r.Table]; !found {
			pgInfo[r.Table] = []string{r.ColumnName}
		} else {
			pgInfo[r.Table] = append(pgInfo[r.Table], r.ColumnName)
		}
	}
	return nil
}
