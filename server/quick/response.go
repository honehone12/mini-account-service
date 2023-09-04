package quick

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func BadRequest() error {
	return echo.NewHTTPError(
		http.StatusBadRequest,
		"input value is not valid",
	)
}

func ServiceError() error {
	return echo.NewHTTPError(
		http.StatusInternalServerError,
		"service has unexpected error",
	)
}

func NotAllowed() error {
	return echo.NewHTTPError(
		http.StatusForbidden,
		"not allowed",
	)
}
