package category_controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/m-row/finder-example/models/category"
)

func (c *Controllers) Store(ctx echo.Context) error {
	var result category.Model
	if err := ctx.Bind(&result); err != nil {
		return err
	}
	// Start transacting
	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if err := c.Models.Category.CreateOne(&result, tx); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return ctx.JSON(http.StatusCreated, result)
}
