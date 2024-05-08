package app

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rwirdemann/datafrog/app/domain"
	"github.com/rwirdemann/datafrog/config"
	"github.com/rwirdemann/datafrog/web/templates"
	log "github.com/sirupsen/logrus"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"time"
)

var client = &http.Client{Timeout: 10 * time.Second}
var Conf config.Config
var apiBaseURL string

func init() {
	Conf = config.NewConfig("config.json")
	apiBaseURL = fmt.Sprintf("http://localhost:%d", Conf.Api.Port)
}

// RegisterHandler registers all known URLs and maps them to their associated
// handlers.
func RegisterHandler(router *mux.Router) {

	// home
	router.HandleFunc("/", IndexHandler)

	// show new form
	router.HandleFunc("/new", NewHandler)

	// start recording
	router.HandleFunc("/create", StartRecording)

	// stop recording
	router.HandleFunc("/stoprecording", StopRecording)

	// delete test
	router.HandleFunc("/delete", DeleteHandler)

	// start verify
	router.HandleFunc("/run", StartVerifyHandler)

	// stop verify
	router.HandleFunc("/stop", StopVerifyHandler)

	// show verification progress
	router.HandleFunc("/verify", ShowVerificationProgressHandler)

	// show test
	router.HandleFunc("/show", ShowHandler)

	// verification progress handler
	router.HandleFunc("/progress", ProgressVerificationHandler)

	// recording progress handler
	router.HandleFunc("/progress-recording", ProgressRecordingHandler)

	// remove expectation from test
	router.HandleFunc("/remove-expectation", RemoveExpectationHandler)
}

func IndexHandler(w http.ResponseWriter, _ *http.Request) {
	allTests := struct {
		Tests []domain.Testcase `json:"tests"`
	}{}
	if r, err := client.Get(fmt.Sprintf("%s/tests", apiBaseURL)); err != nil {
		MsgError = err.Error()
	} else {
		if err := json.NewDecoder(r.Body).Decode(&allTests); err != nil {
			log.Errorf("Error decoding response: %v", err)
		}
	}

	m, e := ClearMessages()
	Render("index.html", w, struct {
		ViewData
		Tests  []domain.Testcase
		Config config.Config
	}{ViewData: ViewData{
		Message: m,
		Title:   "Home",
		Error:   e,
	}, Tests: allTests.Tests, Config: Conf})
}

func ShowHandler(w http.ResponseWriter, r *http.Request) {
	testname := r.URL.Query().Get("testname")
	tc, err := getTestcase(testname, "verifier")
	if err != nil {
		RedirectE(w, r, "/", err)
		return
	}
	m, e := ClearMessages()
	Render("show.html", w, struct {
		ViewData
		Testcase domain.Testcase
	}{ViewData: ViewData{
		Title:   "Show",
		Message: m,
		Error:   e,
	}, Testcase: tc})
}

func ShowVerificationProgressHandler(w http.ResponseWriter, request *http.Request) {
	tc, err := getTestcase(request.FormValue("testname"), "verifier")
	if err != nil {
		RedirectE(w, request, "/", err)
		return
	}

	m, e := ClearMessages()
	Render("verify.html", w, struct {
		ViewData
		Testname     string
		Expectations int
	}{ViewData: ViewData{
		Title:   "Verify",
		Message: m,
		Error:   e,
	}, Testname: tc.Name, Expectations: len(tc.Expectations)})
}

// getTestcase fetches and returns test "name" from the verifier or recorder.
func getTestcase(name string, from string) (domain.Testcase, error) {
	url := fmt.Sprintf("%s/tests/%s?from=%s", apiBaseURL, name, from)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("Error creating request: %v", err)
		return domain.Testcase{}, err
	}
	response, err := client.Do(req)
	if err != nil {
		log.Errorf("Error executing request: %v", err)
		return domain.Testcase{}, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Errorf("Error reading response: %v", err)
		return domain.Testcase{}, err
	}
	var tc domain.Testcase
	if err := json.Unmarshal(body, &tc); err != nil {
		log.Errorf("Error decoding response: %v", err)
		return domain.Testcase{}, err
	}
	return tc, nil
}

// ProgressVerificationHandler renders the partial progress-verification.html
// that shows the progress of the current verification run.
func ProgressVerificationHandler(w http.ResponseWriter, r *http.Request) {
	testname := r.URL.Query().Get("testname")
	progress, err := strconv.Atoi(r.URL.Query().Get("progress"))
	if err != nil {
		return
	}
	tc, err := getTestcase(testname, "verifier")
	if err != nil {
		return
	}
	fulfilled := tc.Fulfilled()
	p := float64(fulfilled) / float64(len(tc.Expectations)) * 100.0
	t, err := template.ParseFS(templates.Templates, "progress-verification.html")
	if err != nil {
		RedirectE(w, r, "/", err)
	}

	color := "is-warning"
	if fulfilled == len(tc.Expectations) {
		color = "is-success"
		progress = 100
	} else {
		progress = int(p)
	}

	data := struct {
		Progress     int
		Testname     string
		Color        string
		Expectations int
		Fulfilled    int
	}{Progress: progress, Testname: testname, Color: color, Expectations: len(tc.Expectations), Fulfilled: fulfilled}
	_ = t.Execute(w, data)
}

// ProgressRecordingHandler renders the partial progress-recording.html
// that shows the progress of the current recording run.
func ProgressRecordingHandler(w http.ResponseWriter, r *http.Request) {
	testname := r.URL.Query().Get("testname")
	progress, err := strconv.Atoi(r.URL.Query().Get("progress"))
	if err != nil {
		return
	}
	tc, err := getTestcase(testname, "recorder")
	if err != nil {
		return
	}
	t, err := template.ParseFS(templates.Templates, "progress-recording.html")
	if err != nil {
		RedirectE(w, r, "/", err)
	}
	progress = len(tc.Expectations) * 3
	data := struct {
		Progress     int
		Testname     string
		Expectations int
	}{Progress: progress, Testname: testname, Expectations: len(tc.Expectations)}
	_ = t.Execute(w, data)
}

// NewHandler renders the new templates
func NewHandler(w http.ResponseWriter, _ *http.Request) {
	RenderS("new.html", w, "New")
}

// StartRecording creates / overrides the test form["testname"] and starts its
// recording.
func StartRecording(w http.ResponseWriter, request *http.Request) {
	testname, err := FormValue(request, "testname")
	if err != nil {
		RedirectE(w, request, "/", err)
		return
	}
	if err := Post(fmt.Sprintf("%s/tests/%s/recordings", apiBaseURL, testname)); err != nil {
		RedirectE(w, request, "/", err)
		return
	}
	MsgSuccess = "Recording has been started. Run UI interactions and click 'Stop recording...' when finished"
	m, e := ClearMessages()
	Render("record.html", w, struct {
		ViewData
		Testname string
	}{ViewData: ViewData{
		Title:   "Record",
		Message: m,
		Error:   e,
	}, Testname: testname})
}

func StopRecording(w http.ResponseWriter, request *http.Request) {
	testname := request.URL.Query().Get("testname")
	url := fmt.Sprintf("%s/tests/%s/recordings", apiBaseURL, testname)
	r, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		RedirectE(w, request, "/", err)
		return
	}
	_, err = client.Do(r)
	if err != nil {
		RedirectE(w, request, "/", err)
		return
	}

	http.Redirect(w, request, "/", http.StatusSeeOther)
}

func DeleteHandler(w http.ResponseWriter, request *http.Request) {
	testname := request.URL.Query().Get("testname")
	url := fmt.Sprintf("%s/tests/%s", apiBaseURL, testname)
	r, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		MsgError = err.Error()
		http.Redirect(w, request, "/", http.StatusSeeOther)
		return
	}
	_, err = client.Do(r)
	if err != nil {
		MsgError = err.Error()
		http.Redirect(w, request, "/", http.StatusSeeOther)
		return
	}
	MsgSuccess = fmt.Sprintf("Test '%s' successfully deleted", testname)
	http.Redirect(w, request, fmt.Sprintf("/"), http.StatusSeeOther)
}

func StartVerifyHandler(w http.ResponseWriter, request *http.Request) {
	testname := request.URL.Query().Get("testname")
	url := fmt.Sprintf("%s/tests/%s/verifications", apiBaseURL, testname)
	r, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		MsgError = err.Error()
		http.Redirect(w, request, "/", http.StatusSeeOther)
		return
	}
	response, err := client.Do(r)
	if err != nil {
		MsgError = err.Error()
		http.Redirect(w, request, "/", http.StatusSeeOther)
		return
	}
	statusOK := response.StatusCode >= 200 && response.StatusCode < 300
	if !statusOK {
		body, _ := io.ReadAll(response.Body)
		MsgError = fmt.Sprintf("HTTP Status: %d => %s", response.StatusCode, body)
	} else {
		MsgSuccess = fmt.Sprintf("Test '%s' has been started. Run test script and click 'Stop...' when you are done!", testname)
	}
	http.Redirect(w, request, fmt.Sprintf("/verify?testname=%s", testname), http.StatusSeeOther)
}

func StopVerifyHandler(w http.ResponseWriter, request *http.Request) {
	testname := request.URL.Query().Get("testname")
	url := fmt.Sprintf("%s/tests/%s/verifications", apiBaseURL, testname)
	r, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		RedirectE(w, request, "/", err)
		return
	}
	response, err := client.Do(r)
	if err != nil {
		RedirectE(w, request, "/", err)
		return
	}

	statusOK := response.StatusCode >= 200 && response.StatusCode < 300
	if !statusOK {
		body, _ := io.ReadAll(response.Body)
		MsgError = fmt.Sprintf("HTTP Status: %d => %s", response.StatusCode, body)
		http.Redirect(w, request, "/", http.StatusSeeOther)
		return
	}
	http.Redirect(w, request, "/show?testname="+testname, http.StatusSeeOther)
}

func RemoveExpectationHandler(http.ResponseWriter, *http.Request) {
}
