package http

import (
	"time"
	"net/http"
	"github.com/labstack/echo"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

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

func HelloWS(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	response := HandlerResponse{"Test websocket endpoint!!!", []*echo.Route{}}

	for {
		err := ws.WriteJSON(response)
		if err != nil {
			return nil
		}

		time.Sleep(time.Second)
	}
}

