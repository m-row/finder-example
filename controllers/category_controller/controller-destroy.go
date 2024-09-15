package category_controller

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/m-row/finder-example/models/category"
)

func (c *Controllers) Destroy(ctx echo.Context) error {
	var result category.Model
	if err := c.ReadUUIDParam(&result.ID, ctx); err != nil {
		return err
	}
	if err := c.Models.Category.GetOne(&result, nil); err != nil {
		return err
	}
	hasChildren, err := c.Models.Category.HasChildren(&result)
	if err != nil {
		return err
	}
	// v, err := c.GetValidator(ctx, result.ModelName())
	// if err != nil {
	// 	return err
	// }
	if hasChildren {
		err := errors.New("can't delete category that has children")
		return err
	}

	// Start transacting
	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if err := c.Models.Category.DeleteOne(&result, nil, tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, result)
}
