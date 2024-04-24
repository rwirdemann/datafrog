package main

import (
	"encoding/json"
	"fmt"
	"github.com/rwirdemann/databasedragon/adapter"
	"github.com/rwirdemann/databasedragon/cmd"
	"github.com/rwirdemann/databasedragon/config"
	"github.com/rwirdemann/databasedragon/httpx/api"
	"github.com/rwirdemann/databasedragon/matcher"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

var verifier *cmd.Verifier
var doneChannels map[string]chan struct{}
var stoppedChannels map[string]chan struct{}

func main() {
	doneChannels = make(map[string]chan struct{})
	stoppedChannels = make(map[string]chan struct{})

	router := mux.NewRouter()
	router.HandleFunc("/tests", AllTests()).Methods("GET")
	router.HandleFunc("/tests/{name}/runs", StartTest()).Methods("PUT")
	router.HandleFunc("/tests/{name}/runs", StopTest()).Methods("DELETE")
	log.Printf("starting http service on port %d...", 3000)
	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, _ := route.GetPathTemplate()
		met, _ := route.GetMethods()
		log.Println(tpl, met)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	err = http.ListenAndServe(fmt.Sprintf(":%d", 3000), router)
	if err != nil {
		log.Fatal(err)
	}
}

func AllTests() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allTests := struct {
			Tests []api.Test `json:"tests"`
		}{}

		entries, err := os.ReadDir(".")
		var tests []string
		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), ".json") && !strings.HasPrefix(entry.Name(), "config") {
				log.Printf("file: %s", entry.Name())
				tests = append(tests, strings.Split(entry.Name(), ".")[0])
			}
		}

		for _, t := range tests {
			var running = false
			if doneChannels[t] != nil {
				running = true
			}
			allTests.Tests = append(allTests.Tests, api.Test{
				Name:    t,
				Running: running,
			})
		}

		b, err := json.Marshal(allTests)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
}

func StartTest() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if len(mux.Vars(request)["name"]) == 0 {
			http.Error(writer, "name is required", http.StatusBadRequest)
			return
		}

		testname := fmt.Sprintf("%s.json", mux.Vars(request)["name"])
		expectationSource, err := adapter.NewFileExpectationSource(testname)
		if err != nil {
			http.Error(writer, "testfile not found", http.StatusNotFound)
			return
		}

		c := config.NewConfig("config.json")
		databaseLog := adapter.NewMYSQLLog(c.Filename)
		t := &adapter.UTCTimer{}
		verifier = cmd.NewVerifier(c, matcher.MySQLTokenizer{}, databaseLog, expectationSource, t, testname)
		doneChannels[testname] = make(chan struct{})
		stoppedChannels[testname] = make(chan struct{})
		go verifier.Start(doneChannels[testname], stoppedChannels[testname])
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.WriteHeader(http.StatusAccepted)
	}
}

func StopTest() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered:", r)
				writer.WriteHeader(http.StatusNotFound)
				return
			}
		}()

		testname := fmt.Sprintf("%s.json", mux.Vars(request)["name"])
		close(doneChannels[testname])
		doneChannels[testname] = nil
		<-stoppedChannels[testname]
		stoppedChannels[testname] = nil
		report := verifier.ReportResults()

		b, err := json.Marshal(report)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		writer.Header().Set("Content-Type", "application/json")
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.Write(b)
	}
}
