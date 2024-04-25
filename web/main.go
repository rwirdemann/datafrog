package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rwirdemann/databasedragon/cmd"
	"github.com/rwirdemann/databasedragon/httpx/api"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"
)

//go:embed templates
var templates embed.FS

var client = &http.Client{Timeout: 10 * time.Second}
var msgSuccess, msgError string

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/new", newHandler)
	router.HandleFunc("/delete", deleteHandler)

	// start recording
	router.HandleFunc("/record", startRecording)

	// stop recording
	router.HandleFunc("/stoprecording", stopRecording)

	router.HandleFunc("/create", createHandler)
	router.HandleFunc("/run", startHandler)
	router.HandleFunc("/stop", stopHandler)
	router.HandleFunc("/show", showHandler)
	log.Println("Listening on :8081...")
	_ = http.ListenAndServe(":8081", router)
}

func newHandler(w http.ResponseWriter, request *http.Request) {
	newtest, err := template.ParseFS(templates, "templates/new.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	newtest.Execute(w, struct {
	}{})
}

func startRecording(w http.ResponseWriter, request *http.Request) {
	record, err := template.ParseFS(templates, "templates/record.html", "templates/header.html", "templates/messages.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	request.ParseForm()
	testname := request.FormValue("testname")
	url := fmt.Sprintf("http://localhost:3000/tests/%s/recordings", testname)
	r, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = client.Do(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	msgSuccess = "Recording has been started. Run UI interactions and click 'Stop recording...' when finished"
	m, e := clearMessages()
	record.Execute(w, struct {
		Testname string
		Message  string
		Error    string
	}{Testname: testname, Message: m, Error: e})
}

func deleteHandler(w http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	testname := request.FormValue("testname")
	url := fmt.Sprintf("http://localhost:3000/tests/%s", testname)
	r, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		http.Redirect(w, request, fmt.Sprintf("/?error=%s", "Error deleting test"), http.StatusSeeOther)
		return
	}
	_, err = client.Do(r)
	if err != nil {
		http.Redirect(w, request, fmt.Sprintf("/?error=%s", "Error deleting test"), http.StatusSeeOther)
		return
	}
	msgSuccess = "Test successfully deleted"
	http.Redirect(w, request, fmt.Sprintf("/"), http.StatusSeeOther)
}

func showHandler(w http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	testname := request.FormValue("testname")
	show, err := template.ParseFS(templates, "templates/show.html", "templates/header.html", "templates/messages.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m, e := clearMessages()
	show.Execute(w, struct {
		Testname string
		Message  string
		Error    string
	}{Testname: testname, Message: m, Error: e})
}

func clearMessages() (string, string) {
	e := msgError
	m := msgSuccess
	msgError = ""
	msgSuccess = ""
	return m, e
}

func startHandler(w http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	testname := request.FormValue("testname")
	url := fmt.Sprintf("http://localhost:3000/tests/%s/runs", testname)
	r, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		msgError = err.Error()
		http.Redirect(w, request, "/", http.StatusSeeOther)
		return
	}
	_, err = client.Do(r)
	if err != nil {
		msgError = err.Error()
		http.Redirect(w, request, "/", http.StatusSeeOther)
		return
	}
	msgSuccess = "Test has been started. Run test script and click 'Stop...' when you are done!"
	http.Redirect(w, request, fmt.Sprintf("/show?testname=%s", testname), http.StatusSeeOther)
}

func stopHandler(w http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	testname := request.FormValue("testname")
	url := fmt.Sprintf("http://localhost:3000/tests/%s/runs", testname)
	r, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response, err := client.Do(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	statusOK := response.StatusCode >= 200 && response.StatusCode < 300
	if !statusOK {
		msgError = fmt.Sprintf("HTTP status not OK: %d", response.StatusCode)
		http.Redirect(w, request, "/", http.StatusSeeOther)
		return
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var report cmd.Report
	err = json.Unmarshal(body, &report)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := template.ParseFS(templates, "templates/result.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result.Execute(w, struct {
		Report cmd.Report
	}{Report: report})
}

func createHandler(w http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	testname := request.FormValue("testname")
	http.Redirect(w, request, fmt.Sprintf("/record&testname=%s", testname), http.StatusSeeOther)
}

func stopRecording(w http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	testname := request.FormValue("testname")
	url := fmt.Sprintf("http://localhost:3000/tests/%s/recordings", testname)
	r, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = client.Do(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, request, "/", http.StatusSeeOther)
}

func indexHandler(w http.ResponseWriter, request *http.Request) {
	r, err := client.Get("http://localhost:3000/tests")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	allTests := struct {
		Tests []api.Test `json:"tests"`
	}{}
	json.NewDecoder(r.Body).Decode(&allTests)

	w.Header().Add("Content-Type", "text/html")
	index, err := template.ParseFS(templates, "templates/index.html", "templates/header.html", "templates/messages.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m, e := clearMessages()
	index.Execute(w, struct {
		Tests   []api.Test
		Message string
		Error   string
	}{Tests: allTests.Tests, Message: m, Error: e})
}
