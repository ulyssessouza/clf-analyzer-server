package data

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	logparser "github.com/Songmu/axslogparser"
)

const MAX_SCORES = 10

const (
	SCORE = iota
	ALERT = iota
)

var Score []SectionScoreEntry
var Alerts []AlertEntry

var db *gorm.DB
var ScoreTicker = time.NewTicker(1 * time.Second)
var AlertTicker = time.NewTicker(1 * time.Second)

// Model for Section's table
type Section struct {
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

// Return type for GetSectionsScore
type AlertEntry struct {
	AlertTime time.Time
	Overcharged bool
}

func OpenDb(dbFilename string) *gorm.DB {
	dbLocal, err := gorm.Open("sqlite3", dbFilename)
	if err != nil {
		panic("failed to connect database")
	}

	return dbLocal
}

func InitDb(dbLocal *gorm.DB) {
	db = dbLocal
	db.AutoMigrate(&Section{})
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
		<-ScoreTicker.C             // Global ScoreTicker triggered every 10s
		*scoreChannel <- SCORE // Trigger endpoint to write the new Score to the client
	}
}

func StartAlertLoop(alertChannel *chan int) {
	for {
		Alerts = GetAlerts(MAX_SCORES)
		<-AlertTicker.C        // Global AlertTicker triggered every 10s
		*alertChannel <- ALERT // Trigger endpoint to write the new Alert to the client
	}
}

func GetSectionsScore(n int) []SectionScoreEntry {
	var sections []SectionScoreEntry
	db.Raw("SELECT COUNT(sections.id) as hits, sections.section FROM sections GROUP BY sections.section ORDER BY COUNT(sections.id) DESC LIMIT ?", n).Scan(&sections)

	return sections
}

func GetAlerts(n int) []AlertEntry {
	var alerts []AlertEntry
	db.Raw("SELECT alerts.created_at as alert_time, alerts.overcharged FROM alerts ORDER BY alerts.created_at DESC LIMIT ?", n).Scan(&alerts)
	return alerts
}

type Count struct {
	N int
}
func CountSectionsIn2Minutes() int {
	var count Count
	now := time.Now()
	last2Minutes := now.Add(-2 * time.Minute) // 2 minutes before
	db.Raw("SELECT COUNT(*) as n FROM sections WHERE sections.created_at > ?", last2Minutes).Scan(&count)
	return count.N
}

func Save(entry interface{}) {
	db.Save(entry)
}