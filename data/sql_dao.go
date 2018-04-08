package data

import (
	"time"
	"github.com/jinzhu/gorm"
)

type SqlDao struct {
	*gorm.DB
}

func NewSqlDao(dbFileName string) *SqlDao {
	sqlDao := &SqlDao{}
	sqlDao.Open(dbFileName)
	return sqlDao
}

// Creates the tables when not there
func (s *SqlDao) Init() {
	s.AutoMigrate(&Log{})
	s.AutoMigrate(&Alert{})
}

func (s *SqlDao) Open(dbFilename string) {
	dbLocal, err := gorm.Open("sqlite3", dbFilename)
	if err != nil {
		panic("failed to connect database")
	}

	s.DB = dbLocal
}

func (s *SqlDao) Close() {
	if s.DB != nil {
		s.Close()
	}
}

func (s *SqlDao) GetSectionsScore(limit int) []SectionScoreEntry {
	var sections []SectionScoreEntry
	s.Raw("SELECT COUNT(logs.id) as hits, logs.section FROM logs GROUP BY logs.section ORDER BY COUNT(logs.id) DESC LIMIT ?", limit).Scan(&sections)
	return sections
}

func (s *SqlDao) GetAlerts(limit int) []Alert {
	var alerts []Alert
	s.Raw("SELECT * FROM alerts ORDER BY alerts.created_at DESC LIMIT ?", limit).Scan(&alerts)
	return alerts
}

// Used by the traffic volume graph
// The values here are fixed because of the presentation layer (termui) sorry :D
func (s *SqlDao) GetAllHitsGroupedBy10Seconds() [120]uint64 {
	var hitEntries [] struct {
		CreatedAt time.Time
	}
	now := time.Now()
	last20Minutes := now.Add(-20 * time.Minute) // 20 minutes of events
	s.Raw("SELECT logs.created_at FROM logs WHERE logs.created_at > ?", last20Minutes).Scan(&hitEntries)

	var ret [120]uint64
	var j, i int64
	for i = 10; i <= 1200; i += 10 {
		for ;int64(len(hitEntries)) > j && hitEntries[j].CreatedAt.Unix() < last20Minutes.Unix() + i; j++ {
			ret[i/10-1]++
		}
	}

	return ret
}

func (s *SqlDao) CountSectionsInDuration(d time.Duration) int {
	var count struct {
		N int
	}
	last2Minutes := time.Now().Add(d) // 2 minutes before
	s.Raw("SELECT COUNT(*) as n FROM logs WHERE logs.created_at > ?", last2Minutes).Scan(&count)
	return count.N
}

func (s *SqlDao) Save(entry interface{}) {
	s.Save(entry)
}