package api

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rwirdemann/datafrog/pkg/df"
	"github.com/rwirdemann/datafrog/pkg/file"
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
	config = df.Config{Filename: "/usr/local/var/mysql/MBP-von-Ralf.log"}
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("/tests/%s/recordings", testname), nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	doneChannel := make(ChannelMap)
	stoppedChannel := make(ChannelMap)
	router := mux.NewRouter()
	router.HandleFunc("/tests/{name}/recordings", StartRecording(doneChannel, stoppedChannel, mocks.LogFactory{})).Methods("POST")
	router.ServeHTTP(rr, req)

	// stop recorded
	close(doneChannel[testFilename])

	// wait until recorded has shut itself down
	<-stoppedChannel[testFilename]

	assert.True(t, file.Exists(testFilename))
	if err := os.Remove(testFilename); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, http.StatusAccepted, rr.Code)
}
