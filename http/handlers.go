package http

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"

	"github.com/ulyssessouza/clf-analyzer-server/data"
)

var upgrader = websocket.Upgrader{}

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
		if !checkAck(ws) {
			break
		}

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
		if !checkAck(ws) {
			break
		}

		alertEntries := getAlertEntriesSlice(data.Alerts)
		err = ws.WriteJSON(alertEntries)
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
		if !checkAck(ws) {
			break
		}

		err := ws.WriteJSON(data.Hits)
		if err != nil {
			return err
		}
		<-hitsChannel // Triggered by data.HitsTicker.C
	}
	return nil
}

func checkAck(ws *websocket.Conn) bool {
	var ack = Ack{}
	err := ws.ReadJSON(ack)
	return err == nil && ack.Code == AckOK
}