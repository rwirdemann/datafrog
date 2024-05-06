package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rwirdemann/datafrog/adapter"
	"github.com/rwirdemann/datafrog/app"
	"github.com/rwirdemann/datafrog/app/domain"
	"github.com/rwirdemann/datafrog/config"
	"github.com/rwirdemann/datafrog/matcher"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var verifier *app.Verifier
var recorder *app.Recorder

// ChannelMap represents a map of empty channels indexed by test names.
type ChannelMap map[string]chan struct{}

// RegisterHandler registers http handler to record and verify testcases.
func RegisterHandler(router *mux.Router,
	verificationDoneChannels ChannelMap, verificationStoppedChannels ChannelMap,
	recordingDoneChannels ChannelMap, recordingStoppedChannels ChannelMap) {

	router.HandleFunc("/tests", AllTests()).Methods("GET")

	// create new test and start recording
	router.HandleFunc("/tests/{name}/recordings", StartRecording(recordingDoneChannels, recordingStoppedChannels)).Methods("POST")

	// stop recording
	router.HandleFunc("/tests/{name}/recordings", StopRecording(recordingDoneChannels, recordingStoppedChannels)).Methods("DELETE")

	// delete test
	router.HandleFunc("/tests/{name}", DeleteTest()).Methods("DELETE")

	// get test
	router.HandleFunc("/tests/{name}", GetTest()).Methods("GET")

	// start verify
	router.HandleFunc("/tests/{name}/verifications", StartVerify(verificationDoneChannels, verificationStoppedChannels)).Methods("PUT")

	// stop verify
	router.HandleFunc("/tests/{name}/verifications", StopVerify(verificationDoneChannels, verificationStoppedChannels)).Methods("DELETE")
}

func GetTest() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		b, err := json.Marshal(verifier.Testcase())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Write(b)
	}
}

// DeleteTest returns a http handler to delete the test given in the request
// param "name".
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

// StartRecording returns a http handler to start the recording of the test
// given in the request params "name". Adds new channels to
// recordingDoneChannels and recordingStoppedChannels. Both channels are used to
// stop the recording (recordingDoneChannels, see StopRecording) and to wait for
// the recorder to gracefully finish its recording loop
// (recordingStoppedChannels see [app.Recorder]).
func StartRecording(recordingDoneChannels ChannelMap, recordingStoppedChannels ChannelMap) http.HandlerFunc {
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
		recorder = app.NewRecorder(c, matcher.MySQLTokenizer{}, databaseLog, recordingSink, t, testname, adapter.GoogleUUIDProvider{})
		recordingDoneChannels[testname] = make(chan struct{})
		recordingStoppedChannels[testname] = make(chan struct{})
		go recorder.Start(recordingDoneChannels[testname], recordingStoppedChannels[testname])
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusAccepted)
	}
}

// StopRecording returns a http handler to stop the recording of the test given
// by the request param "name". Closes the associated channel that is monitored
// by the underlying recording process.
func StopRecording(recordingDoneChannels ChannelMap, recordingStoppedChannels ChannelMap) http.HandlerFunc {
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

		// Notify recorder that recording is done
		close(recordingDoneChannels[testname])
		recordingDoneChannels[testname] = nil

		// Wait until the recorder has gracefully stopped himself
		<-recordingStoppedChannels[testname]
		recordingStoppedChannels[testname] = nil
	}
}

// AllTests returns a http handler that reads all .json files (except
// config.json) in the current directory and returns them as json-encoded http
// response.
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
		_, _ = w.Write(b)
	}
}

func readTestcase(filename string) (domain.Testcase, error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return domain.Testcase{}, err
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			log.Printf("error closing file %s: %v", filename, err)
		}
	}(jsonFile)
	b, _ := io.ReadAll(jsonFile)
	var tc domain.Testcase
	if err := json.Unmarshal(b, &tc); err != nil {
		return domain.Testcase{}, err
	}
	return tc, nil
}

// StartVerify returns a http handler that starts a verification run of the test
// given in the request param "name".
func StartVerify(verificationDoneChannels ChannelMap, verificationStoppedChannels ChannelMap) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
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
		databaseLog := adapter.NewMYSQLLog(c.Filename)
		t := &adapter.UTCTimer{}
		verifier = app.NewVerifier(c, matcher.MySQLTokenizer{}, databaseLog, expectationSource, t, testname)
		verificationDoneChannels[testname] = make(chan struct{})
		verificationStoppedChannels[testname] = make(chan struct{})
		go verifier.Start(verificationDoneChannels[testname], verificationStoppedChannels[testname])
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.WriteHeader(http.StatusAccepted)
	}
}

// StopVerify returns a http handler to stop the verification run of the test
// given in the request param "name".
func StopVerify(verificationDoneChannels ChannelMap, verificationStoppedChannels ChannelMap) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Println("Recovered:", r)
				http.Error(writer, fmt.Sprintf("%v", r), http.StatusInternalServerError)
				return
			}
		}()

		testname := mux.Vars(request)["name"]
		close(verificationDoneChannels[testname])
		verificationDoneChannels[testname] = nil
		<-verificationStoppedChannels[testname]
		verificationStoppedChannels[testname] = nil
		writer.WriteHeader(http.StatusNoContent)
	}
}
