package data

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	logparser "github.com/Songmu/axslogparser"
)

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

var db *gorm.DB

// Real singletons (App members)
var ScoreChannels = NewSynchBroadcastArray(1)
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
}

func OpenDb(dbFilename string) *gorm.DB {
	dbLocal, err := gorm.Open("sqlite3", dbFilename)
	if err != nil {
		panic("failed to connect database")
	}

	return dbLocal
}

// Creates the tables when not there
func InitDb(dbLocal *gorm.DB) {
	db = dbLocal
	db.AutoMigrate(&Log{})
	db.AutoMigrate(&Alert{})
}

func CloseDb() {
	if db != nil {
		db.Close()
	}
}

func StartScoreLoop(scoreChannel *chan int) {
	for {
		Score = GetSectionsScore(MAX_SCORES)
		<-ScoreChannels.C // Global ScoreTicker triggered every 10s
		*scoreChannel <- SCORE // Trigger endpoint to write the new Score to the client
	}
}

func StartAlertLoop(alertChannel *chan int) {
	for {
		Alerts = GetAlerts(MAX_ALERTS)
		<-AlertChannels.C // Global AlertTicker triggered every 10s
		*alertChannel <- ALERT // Trigger endpoint to write the new Alert to the client
	}
}

func StartHitsLoop(hitsChannel *chan int) {
	for {
		Hits = GetAllHitsGroupedBy10Seconds()
		<-HitsChannels.C        // Global HistsTicker triggered every second
		*hitsChannel <- HIT   // Trigger endpoint to write the new Alert to the client
	}
}

func GetSectionsScore(n int) []SectionScoreEntry {
	var sections []SectionScoreEntry
	db.Raw("SELECT COUNT(logs.id) as hits, logs.section FROM logs GROUP BY logs.section ORDER BY COUNT(logs.id) DESC LIMIT ?", n).Scan(&sections)
	return sections
}

func GetAlerts(n int) []Alert {
	var alerts []Alert
	db.Raw("SELECT * FROM alerts ORDER BY alerts.created_at DESC LIMIT ?", n).Scan(&alerts)
	return alerts
}

func GetAllHitsGroupedBy10Seconds() [120]uint64 {
	var hitEntries [] struct {
		CreatedAt time.Time
	}
	now := time.Now()
	last10Minutes := now.Add(-20 * time.Minute) // 20 minutes before
	db.Raw("SELECT logs.created_at FROM logs WHERE logs.created_at > ?", last10Minutes).Scan(&hitEntries)

	var ret [120]uint64
	var j, i int64
	for i = 10; i <= 1200; i += 10 {
		for ;int64(len(hitEntries)) > j && hitEntries[j].CreatedAt.Unix() < last10Minutes.Unix() + i; j++ {
			ret[i/10-1]++
		}
	}

	return ret
}

func CountSectionsIn2Minutes() uint64 {
	var count struct {
		N uint64
	}
	last2Minutes := time.Now().Add(-2 * time.Minute) // 2 minutes before
	db.Raw("SELECT COUNT(*) as n FROM logs WHERE logs.created_at > ?", last2Minutes).Scan(&count)
	return count.N
}

func Save(entry interface{}) {
	db.Save(entry)
}