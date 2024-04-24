package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rwirdemann/databasedragon/httpx/api"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

//go:embed templates
var templates embed.FS

var client = &http.Client{Timeout: 10 * time.Second}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/run", runHandler)
	log.Println("Listening on :8081...")
	_ = http.ListenAndServe(":8081", router)
}

func runHandler(w http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	mode := request.FormValue("mode")
	log.Printf("mode: %s", mode)
	var method string
	if mode == "stop" {
		method = http.MethodDelete
	} else {
		method = http.MethodPut
	}
	url := "http://localhost:3000/tests/create-job/runs"
	r, err := http.NewRequest(method, url, nil)
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
	s := strings.ReplaceAll(string(body), "\n", "%0A")
	indexURL := fmt.Sprintf("/?result=%s", s)
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

	result := request.URL.Query().Get("result")
	var results []string
	if len(strings.Trim(result, " ")) > 0 {
		results = strings.Split(strings.Trim(result, " "), "\n")
	}

	index.Execute(w, struct {
		Title  string
		Tests  []api.Test
		Result []string
	}{Title: "All Tests", Tests: allTests.Tests, Result: results})
}
