package context

import (
	"mini-account-service/db"

	"github.com/labstack/echo/v4"
)

type Metadata struct {
	Name    string
	Version string
}

type Context struct {
	echo.Context
	db.Orm
	Metadata
}
