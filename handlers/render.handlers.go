package handlers

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func renderView(c echo.Context, cmp templ.Component, statusCode ...int) error {
	code := http.StatusOK
	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	c.Response().WriteHeader(code)

	return cmp.Render(c.Request().Context(), c.Response().Writer)
}
