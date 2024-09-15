package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
)

func (app *Application) Serve(e *echo.Echo) error {
	port := os.Getenv("PORT")
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      app.Routes(e, false),
		IdleTimeout:  time.Minute,
		ErrorLog:     log.New(os.Stdout, "", 0),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	log.Printf("starting server at http://localhost%s", srv.Addr)

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		log.Printf("shutting down server, signal: %s", s.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		shutdownError <- srv.Shutdown(ctx)
	}()

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}
	log.Printf("stopped server, Addr: http://localhost%s", srv.Addr)

	return nil
}
