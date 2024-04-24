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

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/new", newHandler)
	router.HandleFunc("/record", recordHandler)
	router.HandleFunc("/stoprecording", stopRecordingHandler)
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

func recordHandler(w http.ResponseWriter, request *http.Request) {
	record, err := template.ParseFS(templates, "templates/record.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	request.ParseForm()
	testname := request.FormValue("testname")
	record.Execute(w, struct {
		Testname string
	}{Testname: testname})
}

func showHandler(w http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	testname := request.FormValue("testname")
	show, err := template.ParseFS(templates, "templates/show.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	show.Execute(w, struct {
		Testname string
	}{Testname: testname})
}

func startHandler(w http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	testname := request.FormValue("testname")
	url := fmt.Sprintf("http://localhost:3000/tests/%s/runs", testname)
	r, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = client.Do(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

func stopRecordingHandler(w http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	testname := request.FormValue("testname")
	log.Printf("testname: %s", testname)
	indexURL := fmt.Sprintf("/?result=%s&recording=false", fmt.Sprintf("Test '%s' has been created.", testname))
	http.Redirect(w, request, indexURL, http.StatusSeeOther)
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
	index, err := template.ParseFS(templates, "templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	index.Execute(w, struct {
		Tests []api.Test
	}{Tests: allTests.Tests})
}
