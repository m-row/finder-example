package api

import (
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GlobalErrorHandler(err error, ctx echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok { //nolint:errorlint // uncomparable
		code = he.Code
	}
	message := "unhandled error"

	switch code {
	case 400:
		message = "bad request"
	case 404:
		message = "not found"
	case 405:
		message = "method not allowed"
	case 409:
		message = "conflict error"
	case 500:
		message = "Internal server error"
		log.Printf("error: %s", err)
	default:
		log.Printf("error: %s", err)
	}
	log.Printf("error: %s", err)

	res := map[string]any{
		"code":    code,
		"message": message,
		"error":   err.Error(),
	}

	if err := ctx.JSON(code, res); err != nil {
		log.Printf("error: %s", err)
	}
}
