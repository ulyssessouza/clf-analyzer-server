package http

import (
	"sync"
	"net/http"

	"github.com/labstack/echo"
	"github.com/gorilla/websocket"
	"github.com/ulyssessouza/clf-analyzer-server/data"
)

var upgrader = websocket.Upgrader{}
var ScoreChannels ScoreChannelsArray

type ScoreChannelsArray struct {
	sync.Mutex
	scoreChannels []chan struct{}
}

func StartListenTicks(c *chan struct{}) {
	for  {
		<-*c
		ScoreChannels.Broadcast()
	}
}

func (s *ScoreChannelsArray) Register(w chan struct{}) {
	s.Lock()
	defer s.Unlock()

	s.scoreChannels = append(s.scoreChannels, w)
}

func (s *ScoreChannelsArray) Deregister(w chan struct{}) {
	s.Lock()
	defer s.Unlock()

	// Delete not including the channel in the new slice
	var newSlice []chan struct{}
	for _, v := range s.scoreChannels {
		if v == w {
			continue
		} else {
			newSlice = append(newSlice, v)
		}
	}

	s.scoreChannels  = newSlice
}

func (s *ScoreChannelsArray) Broadcast() {
	s.Lock()
	defer s.Unlock()

	<-data.Ticker.C
	for _, c := range s.scoreChannels  {
		c <- struct{}{}
	}
}

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

func SectionsScore(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	scoreChannel := make(chan struct{})
	ScoreChannels.Register(scoreChannel)

	defer ws.Close()
	defer ScoreChannels.Deregister(scoreChannel)

	for {
		err := ws.WriteJSON(data.Score)
		if err != nil {
			return err
		}

		<-scoreChannel // Triggered by data.Ticker.C
	}

	return nil
}