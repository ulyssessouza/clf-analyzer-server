package core

import (
	"strings"
	"fmt"
	"time"

	logparser "github.com/Songmu/axslogparser"

	"github.com/ulyssessouza/clf-analyzer-server/data"
)

var AlertHitsThreshold = 10
var ActualHitCount = 0
const twoMinutesAgo = -2 * time.Minute

// Goroutine with the parsing and ingestion loop
func IngestionLoop(saver *data.SaveAndCountInDuration, inputChannel *chan string) {
	for line := range *inputChannel {
		_, log, err := logparser.GuessParser(line)
		if err != nil {
			fmt.Errorf("Ignoring malformed line: %s\n", err)
			continue
		}

		section := data.Log{Log: log, Section: getSection(log.RequestURI)}
		(*saver).Save(&section)
	}
}

// Checks when it should generate a new Alert event or not
// return 1  for overcharged event
// return 0  for still the same state as before
// return -1 for back to normal traffic event
func shouldAlert(actualHitCount int, newHitCount int, hitsThreshold int) int {
	if actualHitCount <= hitsThreshold && newHitCount > hitsThreshold {
		return 1
	} else if actualHitCount > hitsThreshold && newHitCount <= hitsThreshold {
		return -1
	} else {
		return 0
	}
}

// Goroutine that updates alerts list
func UpdateAlertLoop(dao *data.SaveAndCountInDuration){
	for {
		var newHitsCount = (*dao).CountLogsInDuration(twoMinutesAgo)
		switch shouldAlert(ActualHitCount, newHitsCount, AlertHitsThreshold) {
		case  1: (*dao).Save(&data.Alert{Overcharged: true})
		case -1: (*dao).Save(&data.Alert{Overcharged: false})
		}

		ActualHitCount = newHitsCount
		<-time.After(time.Second / 2) // Deliberated time of 1/2 of second :D
	}
}

//From the task description:
//"a section is defined as being what's before the second '/' in a URL. i.e. the section for 'http://my.site.com/pages/create' is 'http://my.site.com/pages'"
//
//Applying this phrase literally. The section of 'http://my.site.com/img.gif' is '/img.gif' and not '/' since the '.' doesn't designate a file in this terms
func getSection(requestURI string) string {
	sections := strings.Split(requestURI, "/")
	section := "/"
	if len(sections) > 0 {
		section = fmt.Sprintf("/%s", sections[1])
	}
	return section
}