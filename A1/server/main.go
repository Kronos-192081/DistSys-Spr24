package main

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/home", func(c echo.Context) error {
		serverNumber := os.Getenv("SERVER_NUMBER")
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Hello from Server: " + serverNumber,
			"status":  "successful",
		})
	})

	e.GET("/heartbeat", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{})
	})

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "5000"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}
