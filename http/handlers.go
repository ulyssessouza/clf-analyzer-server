package http

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"

	"github.com/ulyssessouza/clf-analyzer-server/data"
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

func SectionsScoreHandler(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	scoreChannel := make(chan struct{})
	data.ScoreChannels.Register(scoreChannel)
	defer data.ScoreChannels.Deregister(scoreChannel)

	for {
		err := ws.WriteJSON(data.Score)
		if err != nil {
			return err
		}
		<-scoreChannel // Triggered by data.ScoreTicker.C
	}
	return nil
}

func AlertsHandler(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	alertsChannel := make(chan struct{})
	data.AlertChannels.Register(alertsChannel)
	defer data.AlertChannels.Deregister(alertsChannel)

	for {
		alertEntries := getAlertEntriesSlice(data.Alerts)
		err := ws.WriteJSON(alertEntries)
		if err != nil {
			return err
		}
		<-alertsChannel // Triggered by data.AlertTicker.C
	}
	return nil
}

func HitsHandler(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	hitsChannel := make(chan struct{})
	data.HitsChannels.Register(hitsChannel)
	defer data.HitsChannels.Deregister(hitsChannel)

	for {
		err := ws.WriteJSON(data.Hits)
		if err != nil {
			return err
		}
		<-hitsChannel // Triggered by data.HitsTicker.C
	}
	return nil
}