package core

import (
	"os"
	"testing"

	"github.com/ulyssessouza/clf-analyzer-server/data"
	"github.com/jinzhu/gorm"

	logparser "github.com/Songmu/axslogparser"
)

const dbFilename = "sqlite_testdb.db"
const MAX_RESULTS = 3

func initTestDb() {
	deleteDbFile()
	db, err := gorm.Open("sqlite3", dbFilename)
	db.LogMode(true)
	if err != nil {
		panic("failed to connect database")
	}
	data.InitDb(db)
}

func deleteDbFile() {
	data.CloseDb()

	if _, err := os.Stat(dbFilename); err == nil {
		os.Remove(dbFilename)
	}
}

func TestMain(m *testing.M) {
	initTestDb()
	retCode := m.Run()
	deleteDbFile()
	os.Exit(retCode)
}

func TestGetFirstSections(t *testing.T) {
	const sec1 = "/section1"
	const sec2 = "/section2"
	var sections []data.SectionScoreEntry

	section1 := &data.Section{Log: &logparser.Log{RequestURI: sec1}, Section: sec1}
	data.SaveSection(section1)
	section1 = &data.Section{Log: &logparser.Log{RequestURI: sec1}, Section: sec1}
	data.SaveSection(section1)

	sections = data.GetSectionsScore(MAX_RESULTS)
	if len(sections) != 1 {
		t.Errorf("Got %v expected %v", len(sections), 1)
	}
	if sections[0].Section != sec1 {
		t.Errorf("Got %v expected %v", sections[0].Section, sec1)
	}

	section2 := &data.Section{Log: &logparser.Log{RequestURI: sec2}, Section: sec2}
	data.SaveSection(section2)
	section2 = &data.Section{Log: &logparser.Log{RequestURI: sec2}, Section: sec2}
	data.SaveSection(section2)
	section2 = &data.Section{Log: &logparser.Log{RequestURI: sec2}, Section: sec2}
	data.SaveSection(section2)

	sections = data.GetSectionsScore(MAX_RESULTS)
	if len(sections) != 2 {
		t.Errorf("Got %v expected %v", len(sections), 2)
	}
	if sections[0].Section != sec2 {
		t.Errorf("Got %v expected %v", sections[0].Section, sec2)
	}
}