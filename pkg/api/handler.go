package api

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rwirdemann/datafrog/pkg/df"
	"github.com/rwirdemann/datafrog/pkg/mysql"
	"github.com/rwirdemann/datafrog/pkg/record"
	"github.com/rwirdemann/datafrog/pkg/verify"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var verifier *verify.Verifier
var recorder *record.Recorder
var config df.Config

// ChannelMap represents a map of empty channels indexed by test names.
type ChannelMap map[string]chan struct{}

// RegisterHandler registers http handler to record and verify testcases.
func RegisterHandler(c df.Config, router *mux.Router, verificationDone, verificationStopped, recordingDone, recordingStopped ChannelMap,
	testRepository df.TestRepository) {
	config = c

	// get all tests
	router.HandleFunc("/tests", AllTests(testRepository)).Methods("GET")

	// create new test and start recording
	router.HandleFunc("/tests/{name}/recordings",
		StartRecording(recordingDone, recordingStopped, mysql.LogFactory{}, testRepository)).Methods("POST")

	// stop recording
	router.HandleFunc("/tests/{name}/recordings", StopRecording(recordingDone, recordingStopped)).Methods("DELETE")

	// delete test
	router.HandleFunc("/tests/{name}", DeleteTest()).Methods("DELETE")

	// get test
	router.HandleFunc("/tests/{name}", GetTest(testRepository)).Methods("GET")

	// get recording progress
	router.HandleFunc("/tests/{name}/recordings/progress", GetRecordingProgress()).Methods("GET")

	// get verification progress
	router.HandleFunc("/tests/{name}/verifications/progress", GetVerificationProgress()).Methods("GET")

	// start verify
	router.HandleFunc("/tests/{name}/verifications", StartVerify(verificationDone, verificationStopped)).Methods("PUT")

	// stop verify
	router.HandleFunc("/tests/{name}/verifications", StopVerify(verificationDone, verificationStopped)).Methods("DELETE")

	// channel health
	router.HandleFunc("/channels/{name}/health", ChannelHealth(mysql.LogFactory{})).Methods("GET")
}

func GetRecordingProgress() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tc := recorder.Testcase()
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
		tc := verifier.Testcase()
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

// StartRecording returns a http handler to start the recording of the test given
// in the request params "name". Adds new channels to recordingDoneChannels and
// recordingStoppedChannels. Both channels are used to stop the recording
// (recordingDoneChannels, see StopRecording) and to wait for the recorder to
// gracefully finish its recording loop (recordingStoppedChannels see
// [app.Recorder]).
func StartRecording(done, stopped ChannelMap, logFactory df.LogFactory, repository df.TestRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if len(mux.Vars(r)["name"]) == 0 {
			http.Error(w, "name is required", http.StatusBadRequest)
			return
		}

		if len(config.Channels) == 0 {
			http.Error(w, "at least one channel needs to be configured", http.StatusFailedDependency)
			return
		}

		testname := fmt.Sprintf("%s.json", mux.Vars(r)["name"])
		if repository.Exists(testname) {
			http.Error(w, fmt.Sprintf("test '%s' already exists", testname), http.StatusConflict)
			return
		}

		dbLog := logFactory.Create(config.Channels[0].Log)
		t := &UTCTimer{}
		writer, err := df.NewFileTestWriter(testname)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		recorder = record.NewRecorder(config.Channels[0], mysql.Tokenizer{}, dbLog, writer, t, testname, GoogleUUIDProvider{})
		done[testname] = make(chan struct{})
		stopped[testname] = make(chan struct{})
		go recorder.Start(done[testname], stopped[testname])
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

// StartVerify returns a http handler that starts a verification run of the test
// given in the request param "name".
func StartVerify(verificationDoneChannels ChannelMap, verificationStoppedChannels ChannelMap) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if len(mux.Vars(request)["name"]) == 0 {
			http.Error(writer, "name is required", http.StatusBadRequest)
			return
		}

		if len(config.Channels) == 0 {
			http.Error(writer, "at least one channel needs to be configured", http.StatusFailedDependency)
			return
		}

		// read the test from file
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

		// create a writer to save the test results
		var tw df.TestWriter
		tw, err = df.NewFileTestWriter(fmt.Sprintf("%s.running", testname))
		if err != nil {
			http.Error(writer, err.Error(), http.StatusNotFound)
			return
		}
		databaseLog := mysql.NewMYSQLLog(config.Channels[0].Log)
		t := &UTCTimer{}
		verifier = verify.NewVerifier(config, config.Channels[0], mysql.Tokenizer{}, databaseLog, tc, tw, t, testname)
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
				log.Println("recovered:", r)
				http.Error(writer, fmt.Sprintf("%v", r), http.StatusInternalServerError)
				return
			}
		}()

		testname := mux.Vars(request)["name"]

		// closing done channel forces the verifier to save its testcase
		log.Printf("api: closing done channel")
		close(verificationDoneChannels[testname])

		verificationDoneChannels[testname] = nil

		// wait till verifier has finished its saving
		log.Printf("api: waiting for stopped channel to be closed")
		<-verificationStoppedChannels[testname]
		log.Printf("api: stopped channel closed")

		// copFile .running testfile to original file
		if err := copFile(fmt.Sprintf("%s.running", testname), testname); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		// delete .running file
		if err := deleteFile(fmt.Sprintf("%s.running", testname)); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		verificationStoppedChannels[testname] = nil
		verifier = nil
		writer.WriteHeader(http.StatusNoContent)
	}
}

// ChannelHealth checks the health of the channel "name" by tailing the
// associated log file, triggering the SUT to force a log update and ensures that
// the log file was updated.
func ChannelHealth(lf df.LogFactory) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if len(mux.Vars(request)["name"]) == 0 {
			http.Error(writer, "name is required", http.StatusBadRequest)
			return
		}

		ch, ok := getChannel(mux.Vars(request)["name"])
		if !ok {
			http.Error(writer, "name is required", http.StatusBadRequest)
			return
		}

		clog := lf.Create(ch.Log)

		// jump to logfile end
		err := clog.Tail()
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
		line, err := clog.NextLine(c)
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

var mutex = &sync.Mutex{}

func deleteFile(testname string) error {
	mutex.Lock()
	defer mutex.Unlock()
	log.Printf("api: deleting test file %s", testname)
	return os.Remove(testname)
}

func copFile(src string, dst string) error {
	mutex.Lock()
	defer mutex.Unlock()
	log.Printf("api: copying %s to %s", src, dst)
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func(source *os.File) {
		_ = source.Close()
	}(source)

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func(destination *os.File) {
		_ = destination.Close()
	}(destination)
	_, err = io.Copy(destination, source)
	return err
}
