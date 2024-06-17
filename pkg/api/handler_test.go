package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rwirdemann/datafrog/pkg/df"
	"github.com/rwirdemann/datafrog/pkg/mocks"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestStartRecording(t *testing.T) {
	testname := "create-job"
	testFilename := fmt.Sprintf("%s.json", testname)
	config = df.Config{}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/tests/%s/recordings", testname), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	done := make(ChannelMap)
	stopped := make(ChannelMap)
	r := mux.NewRouter()
	repository := mocks.TestRepository{}
	repository.Testcases = []df.Testcase{}
	r.HandleFunc("/tests/{name}/recordings", StartRecording(done, stopped, mocks.LogFactory{}, repository)).Methods("POST")
	r.ServeHTTP(rr, req)

	// stop recorded
	close(done[testFilename])

	// wait until recorded has shut itself down
	<-stopped[testFilename]

	assert.True(t, exists(testFilename))
	if err := os.Remove(testFilename); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusAccepted, rr.Code)
}

func TestRejectDuplicatedTestname(t *testing.T) {

}

func exists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
