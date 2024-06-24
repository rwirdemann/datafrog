package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rwirdemann/datafrog/pkg/df"
	"github.com/rwirdemann/datafrog/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testname string

func init() {
	testname = "create-job"
}

func TestStartRecordingNoChannels(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/tests/%s/recordings", testname), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	repository := &mocks.TestRepository{}
	r.HandleFunc("/tests/{name}/recordings", StartRecording(mocks.LogFactory{}, repository)).Methods("POST")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusFailedDependency, rr.Code)
}

func TestStartRecording(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/tests/%s/recordings", testname), nil)
	if err != nil {
		t.Fatal(err)
	}
	config.Channels = append(config.Channels, df.Channel{})
	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	repository := &mocks.TestRepository{}
	r.HandleFunc("/tests/{name}/recordings", StartRecording(mocks.LogFactory{}, repository)).Methods("POST")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusAccepted, rr.Code)
	_, ok := runners[testname]
	assert.True(t, ok)
}
