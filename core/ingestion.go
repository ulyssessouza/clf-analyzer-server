package core

import (
	"fmt"

	logparser "github.com/Songmu/axslogparser"

	"github.com/ulyssessouza/clf-analyzer-server/data"
	"strings"
	"time"
)

const AlertShreshold = 10
var ChargeIn2Minutes uint64 = 0

func StartIngestion(inputChannel *chan string) {
	for line := range *inputChannel {
		_, log, err := logparser.GuessParser(line)
		if err != nil {
			fmt.Errorf("Ignoring malformed line: %s\n", err)
			continue
		}

		section := data.Log{Log: log, Section: getSection(log.RequestURI)}
		data.Save(&section)
	}
}

func UpdateAlert() {
	for {
		var countSections = data.CountSectionsIn2Minutes()

		if ChargeIn2Minutes <= AlertShreshold && countSections > AlertShreshold {
			data.Save(&data.Alert{Overcharged: true})
		} else if ChargeIn2Minutes > AlertShreshold && countSections <= AlertShreshold {
			data.Save(&data.Alert{Overcharged: false})
		}

		ChargeIn2Minutes = countSections
		<-time.After(time.Second) // Deliberated time of 1 second :D
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