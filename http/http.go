package http

import (
	"net/http"
	"github.com/labstack/echo"
)

type HandlerResponse struct {
	Message     string
	Endpoints   []*echo.Route
}

// RootHandler godoc
// @Summary List handlers
// @Description lists all the handlers on the app
// @ID root-handler
// @Accept  json
// @Produce  json
// @Success 200 {object} http.HandlerResponse
// @Router / [get]
func RootHandler(c echo.Context) error {
	r := HandlerResponse{"Available endpoints", c.Echo().Routes()}
	return c.JSON(http.StatusOK, r)
}

func TestHandler(c echo.Context) error {
	response := HandlerResponse{"Test endpoint!!!", []*echo.Route{}}
	return c.JSON(http.StatusOK, response)
}

