package http

import (
	"fmt"
	"time"

	"github.com/ulyssessouza/clf-analyzer-server/data"
	"github.com/ulyssessouza/clf-analyzer-server/core"
	"github.com/labstack/echo"
	"github.com/swaggo/echo-swagger"
	"github.com/labstack/echo/middleware"
)

const apiVersion1 = "/v1"

const AckOK = "OK"
type Ack struct {
	Code string
}

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
			data.ScoreChannels.Broadcast()
			break
		case data.ALERT:
			data.AlertChannels.Broadcast()
		case data.HIT:
			data.HitsChannels.Broadcast()
		}
	}
}

func StartHttp(port int) {
	e := echo.New()
	e.HideBanner = true
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/", RootHandler)
	gV1 := e.Group(apiVersion1)
	gV1.GET("/", RootHandler)
	gV1.GET("/swagger/*", echoSwagger.WrapHandler)
	gV1.GET("/score", SectionsScoreHandler)
	gV1.GET("/alert", AlertsHandler)
	gV1.GET("/hits", HitsHandler)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}