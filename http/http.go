package http

import (
	"time"

	"github.com/ulyssessouza/clf-analyzer-server/data"
	"github.com/ulyssessouza/clf-analyzer-server/core"
)

// Real singletons
var scoreChannels data.SynchBroadcastArray
var alertChannels data.SynchBroadcastArray
var hitsChannels data.SynchBroadcastArray

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
		case data.HIT:
			hitsChannels.Broadcast(data.HitsTicker)
		}
	}
}