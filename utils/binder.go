package utils

import (
	"github.com/labstack/echo/v4"
	"time"
)

type CustomBinder struct{}

// Bind method processes form values before binding them to the struct
func (cb *CustomBinder) Bind(i interface{}, c echo.Context) error {
	// Parse the form data first
	if err := c.Request().ParseForm(); err != nil {
		return err
	}

	form := c.Request().Form

	// Convert date strings to time.Time format expected by your struct
	if startDateStr := form.Get("start_date"); startDateStr != "" {
		if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
			form.Set("start_date", t.Format(time.RFC3339))
		}
	}

	if endDateStr := form.Get("end_date"); endDateStr != "" {
		if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
			form.Set("end_date", t.Format(time.RFC3339))
		}
	}

	// Use Echo's default binder to bind the form data to the struct
	binder := new(echo.DefaultBinder)
	return binder.Bind(i, c)
}
