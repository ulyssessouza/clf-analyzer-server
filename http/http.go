package http

import (
	"fmt"
	"log"

	"encoding/json"

	"net/http"

	"github.com/gorilla/mux"
)

type HandlerResponse struct {
	Message     string
	Endpoints   []string
}

var handlersMap = make(map[string] func(http.ResponseWriter, *http.Request))

func WriteResponse(response HandlerResponse, w http.ResponseWriter) {
	js, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)
}

func RootHandler(w http.ResponseWriter, _ *http.Request) {
	var endpoints []string
	for path, _ := range handlersMap {
		endpoints = append(endpoints, path)
	}

	response := HandlerResponse{"Available endpoints", endpoints}
	WriteResponse(response, w)
}

func TestHandler(w http.ResponseWriter, _ *http.Request) {
	response := HandlerResponse{"Test endpoint", []string{}}
	WriteResponse(response, w)
}

func Middleware(toWrap func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		toWrap(w, r)
	}
}

func InitHandlers() {
	handlersMap["/"] = RootHandler
	handlersMap["/test"] = TestHandler
}

func StartHttp() {
	var port = 8000

	var RootRouter = mux.NewRouter()

	InitHandlers()

	for key, value := range handlersMap {
		RootRouter.HandleFunc(key, Middleware(value))
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), RootRouter))
}

