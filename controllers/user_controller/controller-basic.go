package user_controller

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/m-row/finder-example/controller"
	"github.com/m-row/finder-example/models/user"
)

type ControllerBasic struct {
	*controller.Dependencies
}

// Scopes ---------------------------------------------------------------------

// func (c *ControllerBasic) userScope(ctx echo.Context) *uuid.UUID {
// 	scopes := c.Dependencies.CtxScopes(ctx)
// 	if slices.Contains(scopes, "admin") {
// 		return nil
// 	}
// 	return &c.Dependencies.CtxUser(ctx).ID
// }

// Actions --------------------------------------------------------------------

func (c *ControllerBasic) Index(ctx echo.Context) error {
	indexResponse, err := c.Models.User.GetAll(ctx)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, indexResponse)
}

func (c *ControllerBasic) Show(ctx echo.Context) error {
	var result user.Model
	if err := c.ReadUUIDParam(&result.ID, ctx); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	if err := c.Models.User.GetOne(&result); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, result)
}

func (c *ControllerBasic) Store(ctx echo.Context) error {
	result := user.Model{
		CreatedAt: time.Time{},
	}
	// v, err := c.GetValidator(ctx, result.ModelName())
	// if err != nil {
	// 	return err
	// }
	// if valid := result.MergeAndValidate(v); !valid {
	// 	defer v.DeleteNewPicture()
	// 	return c.APIErr.InputValidation(ctx, v)
	// }
	// Start transacting
	tx, err := c.Models.DB.Beginx()
	if err != nil {
		// defer v.DeleteNewPicture()
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	defer func() { _ = tx.Rollback() }()

	if err := c.Models.User.CreateOne(&result, tx); err != nil {
		// defer v.DeleteNewPicture()
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	if err = tx.Commit(); err != nil {
		// defer v.DeleteNewPicture()
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	if err := c.Models.User.GetRoles(&result); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusCreated, result)
}

func (c *ControllerBasic) Update(ctx echo.Context) error {
	var result user.Model
	// v, err := c.GetValidator(ctx, result.ModelName())
	// if err != nil {
	// 	return err
	// }
	if err := c.ReadUUIDParam(&result.ID, ctx); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	if err := c.Models.User.GetOne(&result); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	// if valid := result.MergeAndValidate(v); !valid {
	// 	defer v.DeleteNewPicture()
	// 	return c.APIErr.InputValidation(ctx, v)
	// }
	// Start transacting
	tx, err := c.Models.DB.Beginx()
	if err != nil {
		// defer v.DeleteNewPicture()
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	defer func() { _ = tx.Rollback() }()

	if err := c.Models.User.UpdateOne(&result, tx); err != nil {
		// defer v.DeleteNewPicture()
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ctx.JSON(http.StatusConflict, err.Error())
		default:
			return ctx.JSON(http.StatusInternalServerError, err.Error())
		}
	}

	if err := tx.Commit(); err != nil {
		// defer v.DeleteNewPicture()
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	if err := c.Models.User.GetRoles(&result); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, result)
}

func (c *ControllerBasic) Clear(ctx echo.Context) error {
	var result user.Model
	if err := c.ReadUUIDParam(&result.ID, ctx); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	if err := c.Models.User.GetOne(&result); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	// Start transacting
	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	defer func() { _ = tx.Rollback() }()

	if err := c.Models.User.ClearOne(&result.ID, tx); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}

	if err := tx.Commit(); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	// Delete only if commit succeeds
	// v, err := c.GetValidator(ctx, result.ModelName())
	// if err != nil {
	// 	return err
	// }
	// v.SaveOldImgThumbDists(&result)
	// v.DeleteOldPicture()

	return ctx.JSON(http.StatusOK, result)
}
