package main

import (
	"github.com/labstack/echo/v4"
	"github.com/m-row/finder-example/controller"
	"github.com/m-row/finder-example/controllers"
	"github.com/m-row/finder-example/database"
	"github.com/m-row/finder-example/models"
)

func main() {
	e := echo.New()

	dsn := ""
	db, err := database.OpenSQLX(dsn)
	if err != nil {
		e.Logger.Panicf("couldn't open db: %s", err.Error())
	}
	tc := make(map[string][]string)
	if err := database.PGInfo(db, tc); err != nil {
		e.Logger.Panicf("couldn't get pgInfo: %s", err.Error())
	}

	m := models.Setup(db, tc)
	d := &controller.Dependencies{
		Models: m,
	}
	ctrls := controllers.Setup(d)
	v1 := e.Group("/api/v1")
	rd := &controller.RouterDependencies{
		E: v1,
		// Requires: app.requires,
	}

	ctrls.User.SetAuthRoutes(rd)
	ctrls.User.SetAdminRoutes(rd)
	ctrls.User.SetProfileRoutes(rd)
	ctrls.Role.SetBasicRoutes(rd)
	ctrls.Permission.SetBasicRoutes(rd)

	e.Logger.Fatal(e.Start(":8000"))
}
