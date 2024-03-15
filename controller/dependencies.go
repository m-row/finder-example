package controller

import (
	"errors"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/m-row/finder-example/models"
	"github.com/m-row/finder-example/models/user"
)

// RouterDependencies centralized router grouping dependencies
type RouterDependencies struct {
	// E echo router group usually the /api/v1 group
	E *echo.Group
	// Requires is a middleware that strictly requires a scope to be present
	Requires func(scope ...string) echo.MiddlewareFunc
}

// Dependencies centralized controller dependencies values are injected into
// each controller for easier maintenance
type Dependencies struct {
	// Models application data access layer
	// contains db connection and query builder instance
	// each model represents a database table
	Models *models.Models
}

// ReadUUIDParam parses and validates uuid id parameters.
func (d *Dependencies) ReadUUIDParam(id *uuid.UUID, ctx echo.Context) error {
	paramID := ctx.Param("id")
	parsed, err := uuid.Parse(paramID)
	if err != nil || parsed == uuid.Nil {
		return err
	}
	if id == nil {
		id = &parsed
	}
	*id = parsed
	return nil
}

// ReadIDParam parses and validates integer id parameters.
func (d *Dependencies) ReadIDParam(id *int, ctx echo.Context) error {
	paramID := ctx.Param("id")
	parsed, err := strconv.ParseInt(paramID, 10, 64)
	if err != nil || parsed < 1 {
		return errors.New("invalid id parameter")
	}
	*id = int(parsed)
	return nil
}

func (d *Dependencies) CtxUser(ctx echo.Context) *user.Model {
	ctxUser, ok := ctx.Get("user").(*user.Model)
	if !ok {
		ctxUser = nil
	}
	return ctxUser
}

func (d *Dependencies) CtxScopes(ctx echo.Context) []string {
	ctxScopes, ok := ctx.Get("scopes").([]string)
	if !ok {
		ctxScopes = []string{}
	}
	return ctxScopes
}
