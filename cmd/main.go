package main

import (
	"fmt"

	"github.com/gonext-tech/internal/database"
	"github.com/gonext-tech/internal/handlers"
	"github.com/gonext-tech/internal/routes"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	SECRET_KEY = "SECRET"
	DB_NAME    = "internal"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Error loading .env file")
	}

	e := echo.New()
	e.Static("/", "assets")

	e.HTTPErrorHandler = handlers.CustomHTTPErrorHandler

	// Helpers Middleware
	// e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// Session Middleware
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(SECRET_KEY))))
	e.Use(middleware.Logger())
	db, projectDB, err := database.DBInit()
	if err != nil {
		e.Logger.Fatal(err)
	}
	if err != nil {
		e.Logger.Fatal(err)
	}
	routes.SetupRoutes(e, db, projectDB)

	e.Logger.Fatal(e.Start(":9001"))
}
