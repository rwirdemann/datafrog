package api

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rwirdemann/datafrog/pkg/df"
	"github.com/rwirdemann/datafrog/pkg/file"
	"github.com/rwirdemann/datafrog/pkg/mysql"
	"github.com/rwirdemann/datafrog/pkg/record"
	"github.com/rwirdemann/datafrog/pkg/verify"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var verifier *verify.Verifier
var recorder *record.Recorder
var config df.Config

// ChannelMap represents a map of empty channels indexed by test names.
type ChannelMap map[string]chan struct{}

// RegisterHandler registers http handler to record and verify testcases.
func RegisterHandler(c df.Config, router *mux.Router, verificationDoneChannels ChannelMap, verificationStoppedChannels ChannelMap, recordingDoneChannels ChannelMap, recordingStoppedChannels ChannelMap) {
	config = c
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
	return func(w http.ResponseWriter, r *http.Request) {
		var tc df.Testcase
		from := r.URL.Query().Get("from")
		if from == "verifier" {
			tc = verifier.Testcase()
		} else {
			tc = recorder.Testcase()
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
		databaseLog := mysql.NewMYSQLLog(config.Filename)
		t := &UTCTimer{}

		f, err := os.Create(testname)
		if err != nil {
			log.Fatal(err)
		}
		writer := bufio.NewWriter(f)

		recorder = record.NewRecorder(config, mysql.Tokenizer{}, databaseLog, writer, t, testname, GoogleUUIDProvider{})
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
			Tests []df.Testcase `json:"tests"`
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

func readTestcase(filename string) (df.Testcase, error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return df.Testcase{}, err
	}
	defer func(jsonFile *os.File) {
		err := jsonFile.Close()
		if err != nil {
			log.Printf("error closing file %s: %v", filename, err)
		}
	}(jsonFile)
	b, _ := io.ReadAll(jsonFile)
	var tc df.Testcase
	if err := json.Unmarshal(b, &tc); err != nil {
		return df.Testcase{}, err
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
		expectations, err := os.ReadFile(testname)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusNotFound)
			return
		}
		tc := df.Testcase{}
		if err := json.Unmarshal(expectations, &tc); err != nil {
			http.Error(writer, err.Error(), http.StatusNotFound)
			return
		}

		expectationSource, err := file.NewFileExpectationSource(testname)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusNotFound)
			return
		}

		databaseLog := mysql.NewMYSQLLog(config.Filename)
		t := &UTCTimer{}
		verifier = verify.NewVerifier(config, mysql.Tokenizer{}, databaseLog, tc, expectationSource, t, testname)
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
