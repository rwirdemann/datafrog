package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/rwirdemann/datafrog/pkg/df"
	"github.com/rwirdemann/datafrog/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

var testname string

func init() {
	testname = "create-job"
}

func TestStartRecordingNoChannels(t *testing.T) {
	repository := &mocks.TestRepository{}
	rr := startRecording(t, repository)
	assert.Equal(t, http.StatusFailedDependency, rr.Code)
}

func TestRecording(t *testing.T) {
	config.Channels = append(config.Channels, df.Channel{})
	repository := &mocks.TestRepository{}
	rr := startRecording(t, repository)
	assert.Equal(t, http.StatusAccepted, rr.Code)
	runner, ok := runners[testname]
	assert.True(t, ok)
	runner.Stop()
	tc, err := repository.Get(testname)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "create-job", tc.Name)
}

func TestVerification(t *testing.T) {
	config.Channels = append(config.Channels, df.Channel{})
	repository := &mocks.TestRepository{Testcases: []df.Testcase{{Name: testname}}}
	rr := startVerification(t, repository)
	assert.Equal(t, http.StatusAccepted, rr.Code)
	runner, ok := verifyRunners[testname]
	assert.True(t, ok)
	err := runner.Stop()
	assert.NoError(t, err)
}

func startRecording(t *testing.T, repository df.TestRepository) *httptest.ResponseRecorder {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/tests/%s/recordings", testname), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	r.HandleFunc("/tests/{name}/recordings", StartRecording(repository)).Methods("POST")
	r.ServeHTTP(rr, req)
	return rr
}

func startVerification(t *testing.T, repository df.TestRepository) *httptest.ResponseRecorder {
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/tests/%s/verifications", testname), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	r := mux.NewRouter()
	r.HandleFunc("/tests/{name}/verifications", StartVerification(repository)).Methods("PUT")
	r.ServeHTTP(rr, req)
	return rr
}
