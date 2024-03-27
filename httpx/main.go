package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/tests", CreateTest()).Methods("POST")
	router.HandleFunc("/tests/{test-id}/runs", ResetTest()).Methods("POST")
	router.HandleFunc("/tests/{test-id}/runs", ValidateTest()).Methods("GET")
	log.Printf("starting http service on port %d...", 3000)
	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, _ := route.GetPathTemplate()
		met, _ := route.GetMethods()
		log.Println(tpl, met)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	err = http.ListenAndServe(fmt.Sprintf(":%d", 3000), router)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateTest() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		testID := uuid.New().String()
		writer.Header().Set("Location", testID)
		writer.WriteHeader(http.StatusCreated)
		log.Printf("POST CreateTest: %s\n", testID)
	}
}

func ResetTest() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		testID := mux.Vars(request)["test-id"]
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.WriteHeader(http.StatusOK)
		log.Printf("POST ResetTest: %s\n", testID)
	}
}

func ValidateTest() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		testID := mux.Vars(request)["test-id"]
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.WriteHeader(http.StatusOK)
		log.Printf("POST ValidateTest: %s\n", testID)
	}
}
