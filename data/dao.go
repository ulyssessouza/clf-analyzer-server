package data

import (
	"time"
)

// Interfaces for the data access
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

	GetSuccessesBySection(section string) int
	GetFailsBySection(section string) int
}