package category_controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (c *Controllers) Index(ctx echo.Context) error {
	indexResponse, err := c.Models.Category.GetAll(ctx, nil)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, indexResponse)
}
