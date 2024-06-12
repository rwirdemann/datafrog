package web

import (
	"encoding/json"
	"fmt"
	"github.com/rwirdemann/datafrog/pkg/df"
	"github.com/rwirdemann/datafrog/pkg/driver"
	"github.com/rwirdemann/simpleweb/pkg/simpleweb"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"sort"
	"time"
)

var client *http.Client
var config df.Config
var apiBaseURL string

// Map of channels to synchronize Playwright recording process with data frog web
// app. Map of channels to synchronize Playwright recording process with data
// frog web app. Map of channels to synchronize Playwright recording process with
// data frog web app.
var recordingDoneChannels map[string]chan struct{}

// RegisterHandler registers all known URLs and maps them to their associated
// handlers.
func RegisterHandler(c df.Config) {
	config = c
	recordingDoneChannels = make(map[string]chan struct{})
	client = &http.Client{Timeout: time.Duration(config.Web.Timeout) * time.Second}
	apiBaseURL = fmt.Sprintf("http://localhost:%d", config.Api.Port)

	// home
	simpleweb.Register("/", IndexHandler, "GET")

	// settings
	simpleweb.Register("/settings", SettingsHandler, "GET")

	// show new form
	simpleweb.Register("/new", NewHandler, "GET")

	// start recording
	simpleweb.Register("/create", StartRecording(driver.NewPlaywrightRunner(config)), "POST")

	// stop recording
	simpleweb.Register("/stoprecording", StopRecording, "GET")

	// delete test
	simpleweb.Register("/delete", DeleteHandler, "GET")

	// start test
	simpleweb.Register("/run", StartHandler(driver.NewPlaywrightRunner(config)), "GET")

	// Quit
	simpleweb.Register("/stop", StopHandler, "GET")

	// show test
	simpleweb.Register("/show", ShowHandler, "GET")

	// verification progress handler
	simpleweb.Register("/progress-verification", ProgressVerificationHandler, "GET")

	// recording progress handler
	simpleweb.Register("/progress-recording", ProgressRecordingHandler, "GET")

	// remove expectation from test
	simpleweb.Register("/remove-expectation", RemoveExpectationHandler, "GET")

	simpleweb.Register("/noise", NoiseHandler, "GET")
}

func IndexHandler(w http.ResponseWriter, _ *http.Request) {
	allTests := struct {
		Tests []df.Testcase `json:"tests"`
	}{}
	if r, err := client.Get(fmt.Sprintf("%s/tests", apiBaseURL)); err != nil {
		simpleweb.Error(err.Error())
	} else {
		if err := json.NewDecoder(r.Body).Decode(&allTests); err != nil {
			log.Errorf("Error decoding response: %v", err)
		}
	}

	simpleweb.Render("templates/index.html", w, struct {
		Title string
		Tests []df.Testcase
	}{Title: "Home", Tests: allTests.Tests})
}

func SettingsHandler(w http.ResponseWriter, _ *http.Request) {
	simpleweb.Render("templates/settings.html", w, struct {
		Title  string
		Config df.Config
	}{Title: "Settings", Config: config})
}

func ShowHandler(w http.ResponseWriter, r *http.Request) {
	testname := r.URL.Query().Get("testname")
	url := fmt.Sprintf("%s/tests/%s", apiBaseURL, testname)
	tc, err := getTestcase(url)
	if err != nil {
		simpleweb.RedirectE(w, r, "/", err)
		return
	}
	simpleweb.Render("templates/show.html", w, struct {
		Title    string
		Testcase df.Testcase
	}{Title: "Show", Testcase: tc})
}

// getTestcase fetches and returns test "name".
func getTestcase(url string) (df.Testcase, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("Error creating request: %v", err)
		return df.Testcase{}, err
	}
	response, err := client.Do(req)
	if err != nil {
		log.Errorf("Error executing request: %v", err)
		return df.Testcase{}, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Errorf("Error reading response: %v", err)
		return df.Testcase{}, err
	}
	var tc df.Testcase
	if err := json.Unmarshal(body, &tc); err != nil {
		log.Errorf("Error decoding response: %v", err)
		return df.Testcase{}, err
	}
	return tc, nil
}

// ProgressVerificationHandler renders the partial progress-verification.html
// that shows the progress of the current verification run.
func ProgressVerificationHandler(w http.ResponseWriter, r *http.Request) {
	testname := r.URL.Query().Get("testname")
	url := fmt.Sprintf("%s/tests/%s/verifications/progress", apiBaseURL, testname)
	tc, err := getTestcase(url)
	if err != nil {
		return
	}
	fulfilled := len(tc.Fulfilled())
	p, c := calcProgressAndCssClass(tc)
	if err := simpleweb.RenderPartialE("templates/progress-verification.html", w, struct {
		Progress     int
		Testname     string
		Color        string
		Expectations int
		Fulfilled    int
	}{Progress: p, Testname: testname, Color: c, Expectations: len(tc.Expectations), Fulfilled: fulfilled}); err != nil {
		log.Errorf("Error rendering partial %v", err)
	}
}

func calcProgressAndCssClass(tc df.Testcase) (int, string) {
	fulfilled := len(tc.Fulfilled())
	p := float64(fulfilled) / float64(len(tc.Expectations)) * 100.0
	color := "is-warning"
	progress := int(p)
	if fulfilled == len(tc.Expectations) {
		color = "is-success"
		progress = 100
	}
	return progress, color
}

// ProgressRecordingHandler renders the partial progress-recording.html that
// shows the progress of the current recording run.
func ProgressRecordingHandler(w http.ResponseWriter, r *http.Request) {
	testname := r.URL.Query().Get("testname")
	url := fmt.Sprintf("%s/tests/%s/recordings/progress", apiBaseURL, testname)
	tc, err := getTestcase(url)
	if err != nil {
		log.Errorf("Error getting testcase for progresss from api: %v", err)
		return
	}

	// Determine color and value of progress bar. The bar becomes green and 100% when
	// the recordingDoneChannels[testname] channel was closed. This happens when the
	// user closes the playwright web browser see [driver.PlaywrightRunner.Record].
	var color string
	var progress int
	select {
	default:
		color = "is-warning"
		progress = len(tc.Expectations) * 3
	case <-recordingDoneChannels[testname]:
		color = "is-success"
		progress = 100
	}

	if err := simpleweb.RenderPartialE("templates/progress-recording.html", w, struct {
		Progress     int
		Testname     string
		Color        string
		Expectations int
	}{Progress: progress, Testname: testname, Color: color, Expectations: len(tc.Expectations)}); err != nil {
		log.Errorf("Error rendering partial %v", err)
	}
}

// NewHandler renders the new templates
func NewHandler(w http.ResponseWriter, _ *http.Request) {
	simpleweb.Render("templates/new.html", w, struct {
		Title string
	}{Title: "New Testcase"})
}

// StartRecording creates / overrides the test form["testname"] and starts its
// recording.
func StartRecording(runner driver.PlaywrightRunner) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		testname, err := simpleweb.FormValue(request, "testname")
		if err != nil {
			simpleweb.RedirectE(w, request, "/", err)
			return
		}
		if err := Post(fmt.Sprintf("%s/tests/%s/recordings", apiBaseURL, testname)); err != nil {
			simpleweb.RedirectE(w, request, "/", err)
			return
		}

		if config.AutoVerification {
			recordingDoneChannels[testname] = make(chan struct{})
			go runner.Record(testname, recordingDoneChannels[testname])
		}

		simpleweb.Render("templates/record.html", w, struct {
			Title    string
			Testname string
		}{Title: "Record", Testname: testname})
	}
}

func StopRecording(w http.ResponseWriter, request *http.Request) {
	testname := request.URL.Query().Get("testname")
	url := fmt.Sprintf("%s/tests/%s/recordings", apiBaseURL, testname)
	r, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		simpleweb.RedirectE(w, request, "/", err)
		return
	}
	_, err = client.Do(r)
	if err != nil {
		simpleweb.RedirectE(w, request, "/", err)
		return
	}

	http.Redirect(w, request, fmt.Sprintf("/run?testname=%s.json", testname), http.StatusSeeOther)
}

func DeleteHandler(w http.ResponseWriter, request *http.Request) {
	testname := request.URL.Query().Get("testname")
	url := fmt.Sprintf("%s/tests/%s", apiBaseURL, testname)
	r, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		simpleweb.Error(err.Error())
		http.Redirect(w, request, "/", http.StatusSeeOther)
		return
	}
	_, err = client.Do(r)
	if err != nil {
		simpleweb.Error(err.Error())
		http.Redirect(w, request, "/", http.StatusSeeOther)
		return
	}
	simpleweb.Info(fmt.Sprintf("Test '%s' successfully deleted", testname))
	http.Redirect(w, request, fmt.Sprintf("/"), http.StatusSeeOther)
}

func StartHandler(runner driver.PlaywrightRunner) http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		testname := request.URL.Query().Get("testname")

		// start the test on the api site
		url := fmt.Sprintf("%s/tests/%s/verifications", apiBaseURL, testname)
		r, err := http.NewRequest(http.MethodPut, url, nil)
		if err != nil {
			simpleweb.Error(err.Error())
			http.Redirect(w, request, "/", http.StatusSeeOther)
			return
		}
		response, err := client.Do(r)
		if err != nil {
			simpleweb.Error(err.Error())
			http.Redirect(w, request, "/", http.StatusSeeOther)
			return
		}
		statusOK := response.StatusCode >= 200 && response.StatusCode < 300
		if !statusOK {
			body, _ := io.ReadAll(response.Body)
			simpleweb.Error(fmt.Sprintf("HTTP Status: %d => %s", response.StatusCode, body))
		}

		// start test via driver if configured
		if config.AutoVerification {
			if runner.Exists(testname) {
				simpleweb.Info(fmt.Sprintf("Running auto verfication of test '%s'", testname))
				go runner.Run(testname)
			} else {
				d := fmt.Sprintf("%s/%s", config.Playwright.BaseDir, config.Playwright.TestDir)
				n := runner.ToPlaywright(testname)
				simpleweb.Info(fmt.Sprintf("Playwright test '%s/%s' not found. You have two options:<br>"+
					"1. Create test '%s' in '%s'or <br>2. Run UI interactions manually", d, n, n, d))
			}
		} else {
			simpleweb.Info(fmt.Sprintf("Verification of '%s' started. Run recorded test or execute UI interactions again.", testname))
		}

		// get test progress
		progressUrl := fmt.Sprintf("%s/tests/%s/verifications/progress", apiBaseURL, request.FormValue("testname"))
		tc, err := getTestcase(progressUrl)
		if err != nil {
			simpleweb.RedirectE(w, request, "/", err)
			return
		}

		simpleweb.Render("templates/verify.html", w, struct {
			Title        string
			Testname     string
			Expectations int
		}{Title: fmt.Sprintf("Verifying '%s'...", testname), Testname: tc.Name, Expectations: len(tc.Expectations)})
	}
}

func StopHandler(w http.ResponseWriter, request *http.Request) {
	testname := request.URL.Query().Get("testname")
	url := fmt.Sprintf("%s/tests/%s/verifications", apiBaseURL, testname)
	r, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		simpleweb.RedirectE(w, request, "/", err)
		return
	}
	response, err := client.Do(r)
	if err != nil {
		simpleweb.RedirectE(w, request, "/", err)
		return
	}

	statusOK := response.StatusCode >= 200 && response.StatusCode < 300
	if !statusOK {
		simpleweb.Error("Something went wrong. Please reload page and click on the test to show test results.")
		http.Redirect(w, request, "/", http.StatusSeeOther)
		return
	}
	http.Redirect(w, request, "/show?testname="+testname, http.StatusSeeOther)
}

func NoiseHandler(w http.ResponseWriter, r *http.Request) {
	testname := r.URL.Query().Get("testname")
	url := fmt.Sprintf("%s/tests/%s", apiBaseURL, testname)
	tc, err := getTestcase(url)
	if err != nil {
		simpleweb.RedirectE(w, r, "/", err)
		return
	}
	simpleweb.Render("templates/noise.html", w, struct {
		Title string
		Noise noise
	}{Title: "Noise Sample: " + testname, Noise: buildNoiseData(tc)})
}

func RemoveExpectationHandler(http.ResponseWriter, *http.Request) {
}

type noise struct {
	Verifications int
	EE            []E
}
type E struct {
	Name     string
	Verified int
}

func buildNoiseData(tc df.Testcase) noise {
	noise := noise{Verifications: tc.Verifications}
	for i, e := range tc.Expectations {
		noise.EE = append(noise.EE, E{
			Name:     fmt.Sprintf("E%d: %s", i, e.Shorten(7)),
			Verified: e.Verified,
		})
	}
	sort.Sort(ByVerifications(noise.EE))
	return noise
}

type ByVerifications []E

func (a ByVerifications) Len() int           { return len(a) }
func (a ByVerifications) Less(i, j int) bool { return a[i].Verified > a[j].Verified }
func (a ByVerifications) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
