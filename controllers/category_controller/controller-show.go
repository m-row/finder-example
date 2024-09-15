package category_controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/m-row/finder-example/models/category"
)

func (c *Controllers) Show(ctx echo.Context) error {
	var result category.Model
	if err := c.ReadUUIDParam(&result.ID, ctx); err != nil {
		return err
	}
	if err := c.Models.Category.GetOne(&result, nil); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, result)
}
