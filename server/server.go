package server

import (
	"mini-account-service/db"
	"mini-account-service/server/context"
	"mini-account-service/server/handlers"

	gorillasession "github.com/gorilla/sessions"
	echosession "github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func Run(
	name string,
	version string,
	listenAt string,
	db db.Orm,
	store gorillasession.Store,
) {
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := &context.Context{
				Context: c,
				Orm:     db,
				Metadata: context.Metadata{
					Name:    name,
					Version: version,
				},
			}
			return next(ctx)
		}
	})
	e.Validator = context.NewValidator()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(echosession.Middleware(store))

	e.GET("/", handlers.Root)
	e.POST("/register", handlers.Register)
	e.POST("/login", handlers.Login)
	e.GET("/authorize", handlers.Authorize)

	e.Logger.SetLevel(log.WARN)
	e.Logger.Fatal(e.Start(listenAt))
}
