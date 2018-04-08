package core

import (
	"os"
	"testing"

	"github.com/ulyssessouza/clf-analyzer-server/data"

	logparser "github.com/Songmu/axslogparser"
	"github.com/stretchr/testify/assert"
)

const MAX_RESULTS = 3
const dbFileName = "sqlite_testdb.db"
var sqlDao data.Dao

// Test for the alert logic
func TestShouldAlert(t *testing.T) {
	const alertHitsThreshold = 5
	actualHitCount := alertHitsThreshold
	newHitCount := actualHitCount // Simulate new value for count adding one entry

	assert.Equal(t, 0, shouldAlert(0, newHitCount, alertHitsThreshold))

	actualHitCount = newHitCount
	newHitCount++ // Adding a new log line
	assert.Equal(t, 1, shouldAlert(actualHitCount, newHitCount, alertHitsThreshold))

	actualHitCount = newHitCount
	assert.Equal(t, -1, shouldAlert(actualHitCount, alertHitsThreshold - 1, alertHitsThreshold))
}

func TestGetFirstSections(t *testing.T) {
	const sec1 = "/section1"
	const sec2 = "/section2"
	var sections []data.SectionScoreEntry

	section1 := &data.Log{Log: &logparser.Log{RequestURI: sec1}, Section: sec1}
	sqlDao.Save(section1)
	section1 = &data.Log{Log: &logparser.Log{RequestURI: sec1}, Section: sec1}
	sqlDao.Save(section1)

	sections = sqlDao.GetSectionsScore(MAX_RESULTS)
	if len(sections) != 1 {
		t.Errorf("Got %v expected %v", len(sections), 1)
	}
	if sections[0].Section != sec1 {
		t.Errorf("Got %v expected %v", sections[0].Section, sec1)
	}

	section2 := &data.Log{Log: &logparser.Log{RequestURI: sec2}, Section: sec2}
	sqlDao.Save(section2)
	section2 = &data.Log{Log: &logparser.Log{RequestURI: sec2}, Section: sec2}
	sqlDao.Save(section2)
	section2 = &data.Log{Log: &logparser.Log{RequestURI: sec2}, Section: sec2}
	sqlDao.Save(section2)

	sections = sqlDao.GetSectionsScore(MAX_RESULTS)
	if len(sections) != 2 {
		t.Errorf("Got %v expected %v", len(sections), 2)
	}
	if sections[0].Section != sec2 {
		t.Errorf("Got %v expected %v", sections[0].Section, sec2)
	}
}

func initTestDb() {
	deleteDbFile()

	sqlDao = data.NewSqlDao(dbFileName)
	sqlDao.Init()
}

func deleteDbFile() {
	if sqlDao != nil {
		sqlDao.Close()
	}

	if _, err := os.Stat(dbFileName); err == nil {
		os.Remove(dbFileName)
	}
}

func TestMain(m *testing.M) {
	initTestDb()
	retCode := m.Run()
	deleteDbFile()
	os.Exit(retCode)
}