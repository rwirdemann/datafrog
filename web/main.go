package main

import (
	"embed"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rwirdemann/databasedragon/httpx/api"
	"html/template"
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
	index, err := template.ParseFS(templates, "templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	index.Execute(w, struct {
		Title string
		Tests []api.Test
	}{Title: "All Tests", Tests: allTests.Tests})
}
