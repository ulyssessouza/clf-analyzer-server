package http

import (
	"time"
	"net/http"

	"github.com/labstack/echo"
	"github.com/gorilla/websocket"
	"github.com/ulyssessouza/clf-analyzer-server/data"
	"github.com/ulyssessouza/clf-analyzer-server/core"
)

var upgrader = websocket.Upgrader{}

var scoreChannels data.SynchBroadcastArray
var alertChannels data.SynchBroadcastArray

type AlertEntry struct {
	AlertTime time.Time
	Charge    uint64
	Limit     uint64
}

func getAlertEntriesSlice(alerts []data.Alert) []AlertEntry {
	var alertEntries []AlertEntry
	for _, alert := range alerts {
		var charge uint64
		if alert.Overcharged {
			charge = core.AlertShreshold + 1
		} else {
			charge = 0
		}
		alertEntry := AlertEntry{AlertTime: alert.CreatedAt, Charge: charge, Limit: core.AlertShreshold}
		alertEntries = append(alertEntries, alertEntry)
	}
	return alertEntries
}

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
		alertEntries := getAlertEntriesSlice(data.Alerts)
		err := ws.WriteJSON(alertEntries)
		if err != nil {
			return err
		}
		<-alertsChannel // Triggered by data.AlertTicker.C
	}
	return nil
}