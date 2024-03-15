package role_controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (c *ControllerBasic) GrantAllPermissions(ctx echo.Context) error {
	var id int
	if err := c.ReadIDParam(&id, ctx); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	affected, err := c.Models.Role.GrantAllPermissions(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	response := map[string]any{
		"role_id":    id,
		"grant_type": "all permissions",
		"affected":   affected,
	}
	return ctx.JSON(http.StatusCreated, response)
}

func (c *ControllerBasic) GrantByScope(ctx echo.Context) error {
	var id int
	if err := c.ReadIDParam(&id, ctx); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	if err := ctx.Request().ParseMultipartForm(4 << 20); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	scope := ctx.Request().MultipartForm.Value["scope"][0]
	affected, err := c.Models.Role.GrantByScope(id, scope)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	response := map[string]any{
		"role_id":    id,
		"grant_type": "scope",
		"scope":      scope,
		"affected":   affected,
	}
	return ctx.JSON(http.StatusCreated, response)
}

func (c *ControllerBasic) RevokeAllPermissions(ctx echo.Context) error {
	var id int
	if err := c.ReadIDParam(&id, ctx); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	affected, err := c.Models.Role.RevokeAllPermissions(id)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	response := map[string]any{
		"role_id":     id,
		"revoke_type": "all permissions",
		"affected":    affected,
	}
	return ctx.JSON(http.StatusOK, response)
}
