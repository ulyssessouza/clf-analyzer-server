package data

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	logparser "github.com/Songmu/axslogparser"
)

// UI box size
const MAX_SCORES = 10
const MAX_ALERTS = MAX_SCORES

const (
	SCORE = iota
	ALERT = iota
	HIT = iota
)

var Score []SectionScoreEntry
var Alerts []Alert
var Hits [120]uint64 // Limiting to 120 point in the graph


// Real singletons (App members)
var ScoreChannels = NewSynchBroadcastArray(10)
var AlertChannels = NewSynchBroadcastArray(1)
var HitsChannels = NewSynchBroadcastArray(1)

// Model for Log's table
type Log struct {
	gorm.Model
	*logparser.Log
	Section string
}

// Model for Alert's table
type Alert struct {
	gorm.Model
	Overcharged bool
}

// Return type for GetSectionsScore
type SectionScoreEntry struct {
	Hits uint64
	Section string
	Success int
	Fail int
}

// Goroutine that updates the score list
func StartScoreLoop(dao *Dao, scoreChannel *chan int) {
	for {
		if ScoreChannels.Count() > 0 {
			Score = (*dao).GetSectionsScore(MAX_SCORES)
		}
		<-ScoreChannels.C // Global ScoreTicker triggered every 10s
		*scoreChannel <- SCORE // Trigger endpoint to write the new Score to the client
	}
}

// Goroutine that updates the alert list
func StartAlertLoop(dao *Dao, alertChannel *chan int) {
	for {
		if AlertChannels.Count() > 0 {
			Alerts = (*dao).GetAlerts(MAX_ALERTS)
		}
		<-AlertChannels.C // Global AlertTicker triggered every 10s
		*alertChannel <- ALERT // Trigger endpoint to write the new Alert to the client
	}
}

// Goroutine that updates the hits list
func StartHitsLoop(dao *Dao, hitsChannel *chan int) {
	for {
		if HitsChannels.Count() > 0 {
			Hits = (*dao).GetAllHitsGroupedBy10Seconds()
		}
		<-HitsChannels.C        // Global HistsTicker triggered every second
		*hitsChannel <- HIT   // Trigger endpoint to write the new Alert to the client
	}
}
