package controller

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/m-row/finder-example/models"
)

// RouterDependencies centralized router grouping dependencies
type RouterDependencies struct {
	// E echo router group usually the /api/v1 group
	E *echo.Group
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
