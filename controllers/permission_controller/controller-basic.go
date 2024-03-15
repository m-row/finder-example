package permission_controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/m-row/finder-example/controller"
	"github.com/m-row/finder-example/models/permission"
)

type ControllerBasic struct {
	*controller.Dependencies
}

func (c *ControllerBasic) Index(ctx echo.Context) error {
	indexResponse, err := c.Models.Permission.GetAll(ctx)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, indexResponse)
}

func (c *ControllerBasic) Show(ctx echo.Context) error {
	var result permission.Model
	if err := c.ReadIDParam(&result.ID, ctx); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	if err := c.Models.Permission.GetOne(&result); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, result)
}
