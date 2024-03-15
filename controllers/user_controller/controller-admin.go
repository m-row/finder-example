package user_controller

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/m-row/finder-example/controller"
	"github.com/m-row/finder-example/models/user"
)

type ControllerAdmin struct {
	*controller.Dependencies
}

// Scopes ---------------------------------------------------------------------
// handled by router, all actions are admin

// Actions --------------------------------------------------------------------

func (c *ControllerAdmin) Become(ctx echo.Context) error {
	var result user.Model
	if err := c.ReadUUIDParam(&result.ID, ctx); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	if err := c.Models.User.GetOne(&result); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	tokenResponse, err := result.GenTokenResponse()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	cookie := result.GenCookie(
		tokenResponse.Token,
		time.Now().Add(5*time.Hour),
	)
	ctx.SetCookie(&cookie)
	return ctx.JSON(http.StatusOK, tokenResponse)
}

func (c *ControllerAdmin) GrantRole(ctx echo.Context) error {
	// v, err := c.GetValidator(ctx, result.ModelName())
	// if err != nil {
	// 	return err
	// }

	var roleID int
	var userID uuid.UUID
	// if valid := result.ValidateUserRole(v, &userID, &roleID); !valid {
	// 	return c.APIErr.InputValidation(ctx, v)
	// }

	// Start transacting
	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	defer func() { _ = tx.Rollback() }()
	if err := c.Models.User.GrantRole(&userID, &roleID, tx); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	if err := tx.Commit(); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	responseBody := map[string]string{
		"message": fmt.Sprintf(
			"user %s granted role %d successfully",
			userID,
			roleID,
		),
	}
	return ctx.JSON(http.StatusOK, responseBody)
}

func (c *ControllerAdmin) RevokeRole(ctx echo.Context) error {
	// v, err := c.GetValidator(ctx, result.ModelName())
	// if err != nil {
	// 	return err
	// }

	var roleID int
	var userID uuid.UUID
	// if valid := result.ValidateUserRole(v, &userID, &roleID); !valid {
	// 	return c.APIErr.InputValidation(ctx, v)
	// }
	// Start transacting
	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	defer func() { _ = tx.Rollback() }()
	// update user here
	if err := c.Models.User.RevokeRole(&userID, &roleID, tx); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	// Commit successful transaction
	if err := tx.Commit(); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	responseBody := map[string]string{
		"message": fmt.Sprintf(
			"user %s revoked role %d successfully",
			userID,
			roleID,
		),
	}
	return ctx.JSON(http.StatusOK, responseBody)
}
