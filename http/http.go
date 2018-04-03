package http

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/gorilla/websocket"
	"github.com/ulyssessouza/clf-analyzer-server/data"
)

var upgrader = websocket.Upgrader{}

var scoreChannels data.SynchBroadcastArray
var alertChannels data.SynchBroadcastArray

func StartListenTicks(c *chan int) {
	for {
		signal := <-*c
		switch signal {
		case data.SCORE:
			scoreChannels.Broadcast(data.ScoreTicker)
			break
		case data.ALERT:
			alertChannels.Broadcast(data.AlertTicker)
		}
	}
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

type HandlerResponse struct {
	Message     string
	Endpoints   []*echo.Route
}

func SectionsScore(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	scoreChannel := make(chan struct{})
	scoreChannels.Register(scoreChannel)
	defer scoreChannels.Deregister(scoreChannel)

	for {
		err := ws.WriteJSON(data.Score)
		if err != nil {
			return err
		}
		<-scoreChannel // Triggered by data.ScoreTicker.C
	}
	return nil
}

func Alerts(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	alertsChannel := make(chan struct{})
	alertChannels.Register(alertsChannel)
	defer alertChannels.Deregister(alertsChannel)

	for {
		err := ws.WriteJSON(data.Alerts)
		if err != nil {
			return err
		}
		<-alertsChannel // Triggered by data.AlertTicker.C
	}
	return nil
}