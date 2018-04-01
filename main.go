package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/swaggo/echo-swagger"

	_ "github.com/ulyssessouza/clf-analyzer-server/docs" // docs is generated by Swag CLI
	"github.com/ulyssessouza/clf-analyzer-server/http"
)

func main() {
	port := 8000 // To be set from command line
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", http.RootHandler)
	e.GET("/test", http.TestHandler)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/hello", http.HelloWS)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}
