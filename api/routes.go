package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/m-row/finder-example/controller"
)

func (app *Application) Routes(e *echo.Echo) http.Handler {
	v1 := e.Group("/api/v1")

	deps := &controller.RouterDependencies{
		E: v1,
	}
	app.Controllers.Category.SetRoutes(deps)
	return e
}
