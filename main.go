package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/m-row/finder-example/api"
)

func main() {
	e := echo.New()
	e.HTTPErrorHandler = api.GlobalErrorHandler
	app := api.NewAPI()
	e.Logger.Fatal(app.Serve(e))
}
