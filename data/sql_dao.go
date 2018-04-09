package data

import (
	"time"
	"github.com/jinzhu/gorm"
	"sync"
)

// Type for the Sql data access object
type SqlDao struct {
	sync.RWMutex
	*gorm.DB
}

// Creates a new instance of SqlDao
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

// Opens DB connection
func (s *SqlDao) Open(dbFilename string) {
	dbLocal, err := gorm.Open("sqlite3", dbFilename)
	if err != nil {
		panic("failed to connect database")
	}

	s.DB = dbLocal
}

// Closes DB connection
func (s *SqlDao) Close() {
	if s.DB != nil {
		s.DB.Close()
	}
}

// Gets a certain amount of sections ordered by the most visited
func (s *SqlDao) GetSectionsScore(limit int) []SectionScoreEntry {
	var sections []SectionScoreEntry
	s.RLock()
	s.Raw("SELECT COUNT(logs.id) as hits, logs.section FROM logs GROUP BY logs.section ORDER BY COUNT(logs.id) DESC LIMIT ?", limit).Scan(&sections)
	s.RUnlock()
	return sections
}

func (s *SqlDao) GetAlerts(limit int) []Alert {
	var alerts []Alert
	s.RLock()
	s.Raw("SELECT * FROM alerts ORDER BY alerts.created_at DESC LIMIT ?", limit).Scan(&alerts)
	s.RUnlock()
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
	s.RLock()
	s.Raw("SELECT logs.created_at FROM logs WHERE logs.created_at > ?", last20Minutes).Scan(&hitEntries)
	s.RUnlock()

	var ret [120]uint64
	var j, i int64
	for i = 10; i <= 1200; i += 10 {
		for ;int64(len(hitEntries)) > j && hitEntries[j].CreatedAt.Unix() < last20Minutes.Unix() + i; j++ {
			ret[i/10-1]++
		}
	}

	return ret
}

// Counts logs from a certain time until now
func (s *SqlDao) CountLogsInDuration(d time.Duration) int {
	var count struct {
		N int
	}
	lastRangeOfTime := time.Now().Add(d)
	s.RLock()
	s.Raw("SELECT COUNT(*) as n FROM logs WHERE logs.created_at > ?", lastRangeOfTime).Scan(&count)
	s.RUnlock()
	return count.N
}

func (s *SqlDao) Save(entry interface{}) {
	s.Lock()
	defer s.Unlock()
	s.DB.Save(entry)
}