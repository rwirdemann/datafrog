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

	"github.com/google/uuid"

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

		var running = false
		if doneChannels["create-job.json"] != nil {
			running = true
		}
		allTests.Tests = append(allTests.Tests, api.Test{
			Name:    "create-job",
			Running: running,
		})

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
		c := config.NewConfig("config.json")
		vars := mux.Vars(request)
		testname := fmt.Sprintf("%s.json", vars["name"])
		expectationSource, err := adapter.NewFileExpectationSource(testname)
		if err != nil {
			writer.WriteHeader(http.StatusNotFound)
			return
		}
		databaseLog := adapter.NewMYSQLLog(c.Filename)
		t := &adapter.UTCTimer{}
		verifier = cmd.NewVerifier(c, matcher.MySQLTokenizer{}, databaseLog, expectationSource, t)
		doneChannels[testname] = make(chan struct{})
		stoppedChannels[testname] = make(chan struct{})
		go verifier.Start(doneChannels[testname], stoppedChannels[testname])
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		testID := uuid.New().String()
		writer.Header().Set("Location", testID)
		writer.WriteHeader(http.StatusAccepted)
		log.Printf("PUT StartTest: %s\n", testname)
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
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		close(doneChannels[testname])
		doneChannels[testname] = nil
		<-stoppedChannels[testname]
		stoppedChannels[testname] = nil
		verifier.ReportResults()
		writer.WriteHeader(http.StatusOK)
		log.Printf("DELETE StopTest: %s\n", testname)
	}
}
