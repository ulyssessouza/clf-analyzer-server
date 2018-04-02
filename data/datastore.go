package data

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	logparser "github.com/Songmu/axslogparser"
)

const MAX_SCORES = 10

var db *gorm.DB
var Score []SectionScoreEntry
var Ticker = time.NewTicker(10 * time.Second)

// Model for Section's table
type Section struct {
	gorm.Model
	*logparser.Log
	Section string
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

func InitDb(dbLocal *gorm.DB) {
	db = dbLocal
	db.AutoMigrate(&Section{})
}

func CloseDb() {
	if db != nil {
		db.Close()
	}
}

func GetBySection(sectionId string) []Section {
	var sections []Section
	db.Where("section = ?", sectionId).Find(&sections)
	return sections
}

func StartScoreLoop(scoreChannel *chan struct{}) {
	for {
		Score = GetSectionsScore(MAX_SCORES)
		<-Ticker.C // Global ticker triggered every 10s
		*scoreChannel <- struct{}{} // Trigger endpoint to write the new Score to the client
	}
}

func GetSectionsScore(n int) []SectionScoreEntry {
	var sections []SectionScoreEntry
	db.Raw("SELECT COUNT(sections.id) as hits, sections.section FROM sections GROUP BY sections.section ORDER BY COUNT(sections.id) DESC LIMIT ?", n).Scan(&sections)

	return sections
}

func SaveSection(section *Section) {
	db.Save(section)
}