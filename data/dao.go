package data

import (
	"time"
)

type Dao interface {
	Init()
	Open(dbFilename string)
	Close()

	Save(interface{})

	GetSectionsScore(limit int) []SectionScoreEntry
	GetAlerts(limit int) []Alert
	GetAllHitsGroupedBy10Seconds() [120]uint64
	CountSectionsInDuration(d time.Duration) int
}