package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/m-row/finder-example/api"
	"github.com/m-row/finder-example/controller"
	"github.com/m-row/finder-example/controllers"
	"github.com/m-row/finder-example/database"
	"github.com/m-row/finder-example/models"
)

func main() {
	e := echo.New()
	e.HTTPErrorHandler = api.GlobalErrorHandler

	db, err := database.OpenSQLX()
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
	}

	ctrls.Category.SetRoutes(rd)

	e.Logger.Fatal(e.Start(":8000"))
}
