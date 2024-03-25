package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/tests", CreateTest()).Methods("POST")
	router.HandleFunc("/tests/reset", ResetTest()).Methods("POST")
	router.HandleFunc("/tests/validate", ValidateTest()).Methods("POST")
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
		writer.WriteHeader(http.StatusOK)
		log.Println("POST CreateTest")
	}
}

func ResetTest() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.WriteHeader(http.StatusOK)
		log.Println("POST ResetTest")
	}
}

func ValidateTest() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Access-Control-Allow-Origin", "*")
		writer.WriteHeader(http.StatusOK)
		log.Println("POST ValidateTest")
	}
}
