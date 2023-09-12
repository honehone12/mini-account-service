package context

import (
	"mini-account-service/db"

	"github.com/labstack/echo/v4"
)

type ServiceList struct {
	GamedataService string
}

type Metadata struct {
	Name    string
	Version string
}

type Context struct {
	echo.Context
	db.Orm
	Metadata
	ServiceList
}
