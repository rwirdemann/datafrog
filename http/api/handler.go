package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rwirdemann/databasedragon/adapter"
	"github.com/rwirdemann/databasedragon/app"
	"github.com/rwirdemann/databasedragon/app/domain"
	"github.com/rwirdemann/databasedragon/config"
	"github.com/rwirdemann/databasedragon/matcher"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var verifier *app.Verifier
var doneChannels map[string]chan struct{}
var stoppedChannels map[string]chan struct{}

var recorder *app.Recorder
var recordingDoneChannels map[string]chan struct{}
var recordingStoppedChannels map[string]chan struct{}

func init() {
	doneChannels = make(map[string]chan struct{})
	stoppedChannels = make(map[string]chan struct{})
	recordingDoneChannels = make(map[string]chan struct{})
	recordingStoppedChannels = make(map[string]chan struct{})
}

func RegisterHandler(router *mux.Router) {
	router.HandleFunc("/tests", AllTests()).Methods("GET")

	// create new test and start recording
	router.HandleFunc("/tests/{name}/recordings", CreateTest()).Methods("POST")

	// stop recording
	router.HandleFunc("/tests/{name}/recordings", StopRecording()).Methods("DELETE")

	// delete test
	router.HandleFunc("/tests/{name}", DeleteTest()).Methods("DELETE")

	// start verify
	router.HandleFunc("/tests/{name}/verifications", StartVerify).Methods("PUT")

	// stop verify
	router.HandleFunc("/tests/{name}/verifications", StopVerify()).Methods("DELETE")

}

func DeleteTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(mux.Vars(r)["name"]) == 0 {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}
		testname := mux.Vars(r)["name"]
		err := os.Remove(testname)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func CreateTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(mux.Vars(r)["name"]) == 0 {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}

		testname := fmt.Sprintf("%s.json", mux.Vars(r)["name"])
		c := config.NewConfig("config.json")
		databaseLog := adapter.NewMYSQLLog(c.Filename)
		t := &adapter.UTCTimer{}
		recordingSink := adapter.NewFileRecordingSink(testname)
		recorder = app.NewRecorder(c, matcher.MySQLTokenizer{}, databaseLog, recordingSink, t, testname)
		recordingDoneChannels[testname] = make(chan struct{})
		recordingStoppedChannels[testname] = make(chan struct{})
		go recorder.Start(recordingDoneChannels[testname], recordingStoppedChannels[testname])
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusAccepted)
	}
}

func StopRecording() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered:", r)
				w.WriteHeader(http.StatusNotFound)
				return
			}
		}()

		if len(mux.Vars(r)["name"]) == 0 {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}

		testname := fmt.Sprintf("%s.json", mux.Vars(r)["name"])
		if _, err := os.Stat(testname); os.IsNotExist(err) {
			http.Error(w, "test does not exist", http.StatusNotFound)
			return
		}

		close(recordingDoneChannels[testname])
		recordingDoneChannels[testname] = nil
		<-recordingStoppedChannels[testname]
		recordingStoppedChannels[testname] = nil
	}
}

func AllTests() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allTests := struct {
			Tests []domain.Testcase `json:"tests"`
		}{}

		testfiles, err := os.ReadDir(".")
		for _, f := range testfiles {
			if strings.HasSuffix(f.Name(), ".json") && !strings.HasPrefix(f.Name(), "config") {
				tc, err := readTestcase(f.Name())
				if err != nil {
					log.Printf("error decoding testfile %s: %v", f.Name(), err)
				} else {
					allTests.Tests = append(allTests.Tests, tc)
				}
			}
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

func readTestcase(filename string) (domain.Testcase, error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return domain.Testcase{}, err
	}
	defer jsonFile.Close()
	b, _ := io.ReadAll(jsonFile)
	var tc domain.Testcase
	if err := json.Unmarshal(b, &tc); err != nil {
		return domain.Testcase{}, err
	}
	return tc, nil
}

func StartVerify(writer http.ResponseWriter, request *http.Request) {
	if len(mux.Vars(request)["name"]) == 0 {
		http.Error(writer, "name is required", http.StatusBadRequest)
		return
	}

	testname := mux.Vars(request)["name"]
	expectationSource, err := adapter.NewFileExpectationSource(testname)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusNotFound)
		return
	}

	c := config.NewConfig("config.json")
	databaseLog := adapter.NewMYSQLLog(
		c.Filename)
	t := &adapter.UTCTimer{}
	verifier = app.NewVerifier(c, matcher.MySQLTokenizer{}, databaseLog, expectationSource, t, testname)
	doneChannels[testname] = make(chan struct{})
	stoppedChannels[testname] = make(chan struct{})
	go verifier.Start(doneChannels[testname], stoppedChannels[testname])
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.WriteHeader(http.StatusAccepted)
}

func StopVerify() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered:", r)
				http.Error(writer, fmt.Sprintf("%v", r), http.StatusInternalServerError)
				return
			}
		}()

		testname := mux.Vars(request)["name"]
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
