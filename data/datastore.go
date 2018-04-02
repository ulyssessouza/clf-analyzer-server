package data

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	logparser "github.com/Songmu/axslogparser"
)

var db *gorm.DB

type Section struct {
	gorm.Model
	*logparser.Log
	Section string
}

type SectionScoreEntry struct {
	Hits uint64
	Section string
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

func GetSectionsScore(n int) []SectionScoreEntry {
	var sections []SectionScoreEntry
	db.Raw("SELECT COUNT(sections.id) as hits, sections.section FROM sections GROUP BY sections.section ORDER BY COUNT(sections.id) DESC LIMIT ?", n).Scan(&sections)

	return sections
}

func SaveSection(section *Section) {
	db.Save(section)
}