package main

import (
	"github.com/ulyssessouza/clf-analyzer-server/cmd"
	"github.com/ulyssessouza/clf-analyzer-server/http"
)

func main() {
	cmd.Init()
	http.StartHttp()
}
