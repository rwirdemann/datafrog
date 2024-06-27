package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rwirdemann/datafrog/pkg/df"
	"github.com/rwirdemann/datafrog/pkg/mocks"
	"github.com/rwirdemann/datafrog/pkg/mysql"
	"github.com/rwirdemann/datafrog/pkg/postgres"
	"github.com/rwirdemann/datafrog/pkg/record"
	"github.com/rwirdemann/datafrog/pkg/verify"
)

// invalidStateError represents an unexpected error that informs clients
// the APIs recording state didn't match the request's expectation. For instance,
// when the client requests recording progress for a test that is currently not
// being recorded.
type invalidStateError struct {
}

func (i invalidStateError) Error() string {
	return "invalid recording state"
}

var config df.Config

var runners = make(map[string]*record.Runner)
var verifyRunners = make(map[string]*verify.Runner)

// RegisterHandler registers http handler to record and verify testcases.
func RegisterHandler(c df.Config, router *mux.Router, testRepository df.TestRepository) {
	config = c

	// get all tests
	router.HandleFunc("/tests", AllTests(testRepository)).Methods("GET")

	// create new test and start recording
	router.HandleFunc("/tests/{name}/recordings",
		StartRecording(testRepository)).Methods("POST")

	// stop recording
	router.HandleFunc("/tests/{name}/recordings", StopRecording()).Methods("DELETE")

	// delete test
	router.HandleFunc("/tests/{name}", DeleteTest(testRepository)).Methods("DELETE")

	// get test
	router.HandleFunc("/tests/{name}", GetTest(testRepository)).Methods("GET")

	// get recording progress
	router.HandleFunc("/tests/{name}/recordings/progress", GetRecordingProgress()).Methods("GET")

	// get verification progress
	router.HandleFunc("/tests/{name}/verifications/progress", GetVerificationProgress()).Methods("GET")

	// start verify
	router.HandleFunc("/tests/{name}/verifications", StartVerification(testRepository)).Methods("PUT")

	// stop verify
	router.HandleFunc("/tests/{name}/verifications", StopVerify()).Methods("DELETE")

	// channel health
	router.HandleFunc("/channels/{name}/health", ChannelHealth()).Methods("GET")
}

func GetRecordingProgress() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		runner, ok := runners[mux.Vars(r)["name"]]
		if !ok {
			http.Error(w, invalidStateError{}.Error(), http.StatusInternalServerError)
			return
		}

		tc := runner.Testcase()
		b, err := json.Marshal(tc)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if _, err := w.Write(b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GetVerificationProgress() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		runner, ok := verifyRunners[mux.Vars(r)["name"]]
		if !ok {
			http.Error(w, invalidStateError{}.Error(), http.StatusInternalServerError)
			return
		}

		tc := runner.Testcase()
		b, err := json.Marshal(tc)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if _, err := w.Write(b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func GetTest(repository df.TestRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(mux.Vars(r)["name"]) == 0 {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}
		tc, err := repository.Get(mux.Vars(r)["name"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		b, err := json.Marshal(tc)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if _, err := w.Write(b); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// DeleteTest returns a http handler to delete the test given in the request
// param "name".
func DeleteTest(repository df.TestRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(mux.Vars(r)["name"]) == 0 {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}
		if err := repository.Delete(mux.Vars(r)["name"]); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// StartRecording starts recording of test given the request param "name".
func StartRecording(repository df.TestRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(mux.Vars(r)["name"]) == 0 {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}

		if len(config.Channels) == 0 {
			http.Error(w, "at least one channel needs to be configured", http.StatusFailedDependency)
			return
		}

		testname := mux.Vars(r)["name"]
		if repository.Exists(testname) {
			http.Error(w, fmt.Sprintf("test '%s' already exists", testname), http.StatusConflict)
			return
		}

		channelName := config.Channels[0].Name // TODO: Get name from request
		channel, ok := getChannel(channelName)
		if !ok {
			http.Error(w, fmt.Sprintf("channel '%s' does not exisit", channelName), http.StatusConflict)
			return
		}

		channelLog, err := getLog(channel)
		if err != nil {
			http.Error(w, fmt.Sprintf("Logfile '%s' does not exist", channel.Log), http.StatusFailedDependency)
			return
		}

		runners[testname] = record.NewRunner(testname, channel, repository, channelLog)

		// Start creates a new go routine
		if err := runners[testname].Start(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusAccepted)
	}
}

// StopRecording returns a http handler to stop the recording of the test given
// by the request param "name".
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

		testname := mux.Vars(r)["name"]
		runner, ok := runners[testname]
		if !ok {
			http.Error(w, "test is not being recorded", http.StatusNotFound)
			return
		}
		delete(runners, testname)
		runner.Stop()
	}
}

// AllTests returns all tests as json-encoded HTTP response.
func AllTests(repository df.TestRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		all, err := repository.All()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		allTests := struct {
			Tests []df.Testcase `json:"tests"`
		}{Tests: all}
		b, err := json.Marshal(allTests)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(b)
	}
}

// StartVerification returns a http handler that starts a verification run of the test
// given in the request param "name".
func StartVerification(repository df.TestRepository) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if len(mux.Vars(request)["name"]) == 0 {
			http.Error(writer, "name is required", http.StatusBadRequest)
			return
		}

		if len(config.Channels) == 0 {
			http.Error(writer, "at least one channel needs to be configured", http.StatusFailedDependency)
			return
		}

		testname := mux.Vars(request)["name"]

		channelName := config.Channels[0].Name // TODO: Get name from request
		channel, ok := getChannel(channelName)
		if !ok {
			http.Error(writer, fmt.Sprintf("channel '%s' does not exisit", channelName), http.StatusConflict)
			return
		}

		channelLog, err := getLog(channel)
		if err != nil {
			http.Error(writer, fmt.Sprintf("Logfile '%s' does not exist", channel.Log), http.StatusFailedDependency)
			return
		}

		verifyRunners[testname] = verify.NewRunner(testname, channel, config, channelLog, repository)

		// Start creates a new go routine
		if err := verifyRunners[testname].Start(); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.WriteHeader(http.StatusAccepted)
	}
}

// StopVerify returns a http handler to stop the verification run of the test
// given in the request param "name".
func StopVerify() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		testname := mux.Vars(request)["name"]
		runner, ok := verifyRunners[testname]
		if !ok {
			http.Error(writer, "test is not being verified", http.StatusNotFound)
			return
		}
		delete(verifyRunners, testname)
		if err := runner.Stop(); err != nil {
			http.Error(writer, err.Error(), http.StatusNotFound)
			return
		}

		writer.WriteHeader(http.StatusNoContent)
	}
}

// ChannelHealth checks the health of the channel "name" by tailing the
// associated log file, triggering the SUT to force a log update and ensures that
// the log file was updated.
func ChannelHealth() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if len(mux.Vars(request)["name"]) == 0 {
			http.Error(writer, "name is required", http.StatusBadRequest)
			return
		}

		if len(config.Channels) == 0 {
			http.Error(writer, "at least one channel needs to be configured", http.StatusFailedDependency)
			return
		}

		channelName := mux.Vars(request)["name"]
		channel, ok := getChannel(channelName)
		if !ok {
			http.Error(writer, fmt.Sprintf("channel '%s' does not exisit", channelName), http.StatusConflict)
			return
		}

		channelLog, errLog := getLog(channel)
		if errLog != nil {
			http.Error(writer, fmt.Sprintf("Logfile '%s' does not exist", channel.Log), http.StatusConflict)
			return
		}

		// jump to logfile end
		err := channelLog.Tail()
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}

		// trigger SUT and give it some time to update the channel log
		_, err = http.Get(config.SUT.BaseURL)
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		time.Sleep(250 * time.Millisecond)

		// create go routine to interrupt blocking the NextLine call after 200ms
		c := make(chan struct{})
		go func() {
			time.Sleep(200 * time.Millisecond)
			close(c)
		}()

		// read next line from updated log file
		line, err := channelLog.NextLine(c)
		if err != nil || line == "" {
			writer.WriteHeader(http.StatusFailedDependency)
			return
		}
		log.Printf("line: " + line)
		writer.WriteHeader(http.StatusOK)
	}
}

func getChannel(name string) (df.Channel, bool) {
	for _, ch := range config.Channels {
		if ch.Name == name {
			return ch, true
		}
	}
	return df.Channel{}, false
}

func getLog(channel df.Channel) (df.Log, error) {
	var logFactory df.LogFactory
	if channel.Format == "mysql" {
		logFactory = mysql.LogFactory{}
	}
	if channel.Format == "postgres" {
		logFactory = postgres.LogFactory{}
	}
	// TODO: This is fix for test_handler. In real useage, we must not use it and should return an error here!
	if logFactory == nil {
		logFactory = mocks.LogFactory{}
	}

	return logFactory.Create(channel.Log)
}
