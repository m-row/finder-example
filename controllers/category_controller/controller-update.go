package category_controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/m-row/finder-example/models/category"
)

func (c *Controllers) Update(ctx echo.Context) error {
	var result category.Model
	// v, err := c.GetValidator(ctx, result.ModelName())
	// if err != nil {
	// 	return err
	// }
	if err := c.ReadUUIDParam(&result.ID, ctx); err != nil {
		return err
	}
	if err := c.Models.Category.GetOne(&result, nil); err != nil {
		return err
	}
	// if valid := result.MergeAndValidate(v); !valid {
	// 	defer v.DeleteNewPicture()
	// 	return c.APIErr.InputValidation(ctx, v)
	// }
	// Start transacting
	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	if err := c.Models.Category.UpdateOne(&result, nil, tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, result)
}
