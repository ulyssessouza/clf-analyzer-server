package data

import (
	"time"
)

type SaveAndCountInDuration interface {
	Save(interface{})
	CountLogsInDuration(d time.Duration) int
}

type Dao interface {
	SaveAndCountInDuration

	Init()
	Open(dbFilename string)
	Close()


	GetSectionsScore(limit int) []SectionScoreEntry
	GetAlerts(limit int) []Alert
	GetAllHitsGroupedBy10Seconds() [120]uint64
}