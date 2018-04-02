package core

import (
	"bufio"
	"fmt"
	"os"

	logparser "github.com/Songmu/axslogparser"

	"github.com/ulyssessouza/clf-analyzer-server/data"
	"strings"
	"github.com/jinzhu/gorm"
)

func StartIngestion(f *os.File) {
	db, err := gorm.Open("sqlite3", "sqlite_clf_analyzer.db")
	if err != nil {
		panic("failed to connect database")
	}
	data.InitDb(db)
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

		data.SaveSection(&section)
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