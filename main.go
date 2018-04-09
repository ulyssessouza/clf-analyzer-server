package main

import (
	"flag"

	"github.com/ulyssessouza/clf-analyzer-server/core"
	"github.com/ulyssessouza/clf-analyzer-server/data"
	"github.com/ulyssessouza/clf-analyzer-server/http"
)

func startGoroutines(dao *data.Dao, tailFlag *string, inputLineChan *chan string, cacheRefreshChan *chan int) {
	// Choose input mode
	if tailFlag != nil && *tailFlag != "" && *tailFlag != "stdin" {
		go inputFromTail(inputLineChan, *tailFlag)
	} else {
		go inputFromStdIn(inputLineChan)
	}

	var daoSaver = data.SaveAndCountInDuration(*dao)

	go core.IngestionLoop(&daoSaver, inputLineChan) // Starts ingesting lines included in the channel by the previous goroutines

	go http.StartListenTicks(cacheRefreshChan) // Listen ticks

	go data.StartScoreLoop(dao, cacheRefreshChan)
	go data.StartAlertLoop(dao, cacheRefreshChan)
	go data.StartHitsLoop(dao, cacheRefreshChan)
	go core.UpdateAlertLoop(&daoSaver)
}

func main() {
	var port = flag.Int("port", 8000, "port to listen on")
	var tailFlag = flag.String("tail", "stdin", "file to tail")
	var alertHitsThreshold = flag.Int("alertThreshold", 10, "threshold of log entries in 2 minutes to consider a normal traffic.")
	flag.Parse()

	core.AlertHitsThreshold = *alertHitsThreshold

	var sqlDao data.Dao = data.NewSqlDao("sqlite_clf_analyzer.db")
	sqlDao.Init()
	defer sqlDao.Close()

	var cacheRefreshChan = make(chan int) // data.ticker -> SQL select -> scoreChannels.Broadcast()
	var inputLineChan = make(chan string) // Channel used to make the input source generic
	defer close(inputLineChan)
	defer close(cacheRefreshChan)

	startGoroutines(&sqlDao, tailFlag, &inputLineChan, &cacheRefreshChan)
	http.StartHttp(*port)
}
