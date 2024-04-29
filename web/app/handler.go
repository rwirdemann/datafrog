package app

import (
	"encoding/json"
	"fmt"
	"github.com/rwirdemann/databasedragon/app/domain"
	"github.com/rwirdemann/databasedragon/config"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

var client = &http.Client{Timeout: 10 * time.Second}
var Conf config.Config
var apiBaseURL string

func init() {
	Conf = config.NewConfig("config.json")
	apiBaseURL = fmt.Sprintf("http://localhost:%d", Conf.Api.Port)
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

func ShowHandler(w http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	m, e := ClearMessages()
	Render("show.html", w, struct {
		ViewData
		Testname string
	}{ViewData: ViewData{
		Title:   "Show",
		Message: m,
		Error:   e,
	}, Testname: request.FormValue("testname")})
}

func NewHandler(w http.ResponseWriter, _ *http.Request) {
	RenderS("new.html", w, "New")
}

func StartRecording(w http.ResponseWriter, request *http.Request) {
	if err := request.ParseForm(); err != nil {
		RedirectE(w, request, "/", err)
		return
	}
	testname := request.FormValue("testname")
	url := fmt.Sprintf("%s/tests/%s/recordings", apiBaseURL, testname)
	r, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		RedirectE(w, request, "/", err)
		return
	}
	_, err = client.Do(r)
	if err != nil {
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

func StartHandler(w http.ResponseWriter, request *http.Request) {
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
	http.Redirect(w, request, fmt.Sprintf("/show?testname=%s", testname), http.StatusSeeOther)
}

func StopHandler(w http.ResponseWriter, request *http.Request) {
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

	body, err := io.ReadAll(response.Body)
	if err != nil {
		RedirectE(w, request, "/", err)
		return
	}
	var report domain.Report
	err = json.Unmarshal(body, &report)
	if err != nil {
		RedirectE(w, request, "/", err)
		return
	}

	m, e := ClearMessages()
	Render("result.html", w, struct {
		ViewData
		Testname string
		Report   domain.Report
	}{ViewData: ViewData{
		Title:   "Result",
		Message: m,
		Error:   e,
	}, Testname: testname, Report: report})
}

func CreateHandler(w http.ResponseWriter, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		RedirectE(w, request, "/", err)
		return
	}
	http.Redirect(w, request, fmt.Sprintf("/record&testname=%s", request.FormValue("testname")), http.StatusSeeOther)
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
