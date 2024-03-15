package role_controller

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/m-row/finder-example/controller"
	"github.com/m-row/finder-example/models/role"
)

type ControllerBasic struct {
	*controller.Dependencies
}

func (c *ControllerBasic) Index(ctx echo.Context) error {
	indexResponse, err := c.Models.Role.GetAll(ctx)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, indexResponse)
}

func (c *ControllerBasic) Show(ctx echo.Context) error {
	var result role.Model
	if err := c.ReadIDParam(&result.ID, ctx); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	if err := c.Models.Role.GetOne(&result); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, result)
}

func (c *ControllerBasic) Store(ctx echo.Context) error {
	var result role.Model
	// v, err := c.GetValidator(ctx, result.ModelName())
	// if err != nil {
	// 	return err
	// }
	// if valid := result.MergeAndValidate(v); !valid {
	// 	return c.APIErr.InputValidation(ctx, v)
	// }
	// Start transacting
	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	defer func() { _ = tx.Rollback() }()

	if err := c.Models.Role.CreateOne(&result, tx); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	if err = tx.Commit(); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusCreated, result)
}

func (c *ControllerBasic) Update(ctx echo.Context) error {
	var result role.Model
	// v, err := c.GetValidator(ctx, result.ModelName())
	// if err != nil {
	// 	return err
	// }
	if err := c.ReadIDParam(&result.ID, ctx); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	if err := c.Models.Role.GetOne(&result); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	// if valid := result.MergeAndValidate(v); !valid {
	// 	return c.APIErr.InputValidation(ctx, v)
	// }
	// Start transacting
	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	defer func() { _ = tx.Rollback() }()

	if err := c.Models.Role.UpdateOne(&result, tx); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ctx.JSON(http.StatusConflict, err.Error())
		default:
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	if err := tx.Commit(); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, result)
}

func (c *ControllerBasic) Destroy(ctx echo.Context) error {
	var result role.Model
	if err := c.ReadIDParam(&result.ID, ctx); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	if result.ID == 1 {
		err := errors.New("you shouldn't do that")
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	// Start transacting
	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	defer func() { _ = tx.Rollback() }()

	if err := c.Models.Role.DeleteOne(&result, tx); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	if err := tx.Commit(); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, result)
}
