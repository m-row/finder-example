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

type ControllerAuth struct {
	*controller.Dependencies
}

func (c *ControllerAuth) Login(ctx echo.Context) error {
	var result user.Model
	comparePassword := false
	// v, err := c.GetValidator(ctx, result.ModelName())
	// if err != nil {
	// 	return err
	// }
	// valid, err := result.MergeLogin(v, nil, &comparePassword)
	// if err != nil {
	// 	return c.APIErr.Firebase(ctx, err)
	// }
	// if !valid {
	// 	return c.APIErr.InputValidation(ctx, v)
	// }
	if err := c.Models.User.GetOne(&result); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ctx.JSON(http.StatusUnauthorized, err.Error())
		}
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	// does not compare password if firebase_id_token is provided or is
	// registration
	if comparePassword {
		if ok, err := result.Password.
			Match(result.PasswordHash); err != nil || !ok {
			return ctx.JSON(http.StatusUnauthorized, err.Error())
		}
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

func (c *ControllerAuth) Logout(ctx echo.Context) error {
	var model user.Model
	var token user.Token
	var accessToken string
	var message string
	var status int

	reqCookie, err := ctx.Cookie("accessToken")
	if err != nil {
		accessToken = ""
	} else {
		accessToken = reqCookie.Value
	}
	if accessToken != "" {
		message = "logged out"
		status = http.StatusOK
	} else {
		message = "not logged in..."
		status = http.StatusUnauthorized
	}
	cookie := model.GenCookie(
		token,
		time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
	)
	ctx.SetCookie(&cookie)
	ctx.Set("scopes", []string{})
	response := map[string]any{
		"status":     status,
		"message":    message,
		"request_id": ctx.Response().Header().Get(echo.HeaderXRequestID),
		"errors":     nil,
	}
	return ctx.JSON(status, response)
}
