package core

import (
	"bufio"
	"fmt"
	"os"

	logparser "github.com/Songmu/axslogparser"

	"github.com/ulyssessouza/clf-analyzer-server/data"
	"strings"
	"time"
)

const ALERT_THRESHOLD = 10
var overCharged = false

func StartIngestion(f *os.File) {
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()

		parser, log, err := logparser.GuessParser(line)
		if err != nil {
			fmt.Errorf("%s\n", err)
		}
		if _, ok := parser.(*logparser.Apache); !ok {
			fmt.Errorf("Invalid format: %s\n", line)
		}

		section := data.Section{Log: log, Section: getSection(log.RequestURI)}
		data.Save(&section)
	}
}

func UpdateAlert() {
	for {
		countSections := data.CountSectionsIn2Minutes()
		fmt.Printf("count %d\n", countSections)
		if !overCharged && countSections > ALERT_THRESHOLD {
			overCharged = true
			data.Save(&data.Alert{Overcharged: true})
		} else if overCharged && countSections <= ALERT_THRESHOLD {
			overCharged = false
			data.Save(&data.Alert{Overcharged: false})
		}
		<-time.After(time.Second)
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