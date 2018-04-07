package main

import (
	"bufio"
	"fmt"
	"github.com/hpcloud/tail"
	"os"
)

func inputFromTail(c *chan string, fileName string) {
	t, err := tail.TailFile(fileName, tail.Config{Follow: true, ReOpen: true})
	if err != nil {
		panic(fmt.Sprintf("Couldn't tail file '%s'", fileName))
	}

	for line := range t.Lines {
		*c <- line.Text
	}
}

func inputFromStdIn(c *chan string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		*c <- scanner.Text()
	}
}
