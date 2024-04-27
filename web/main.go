package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rwirdemann/databasedragon/web/app"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()

	// home
	router.HandleFunc("/", app.IndexHandler)

	// show new form
	router.HandleFunc("/new", app.NewHandler)

	router.HandleFunc("/create", app.CreateHandler)

	// delete test
	router.HandleFunc("/delete", app.DeleteHandler)

	// start recording
	router.HandleFunc("/record", app.StartRecording)

	// stop recording
	router.HandleFunc("/stoprecording", app.StopRecording)

	// start verifx
	router.HandleFunc("/run", app.StartHandler)

	// stop verify
	router.HandleFunc("/stop", app.StopHandler)

	// show test
	router.HandleFunc("/show", app.ShowHandler)

	log.Printf("Listening on :%d...", app.Conf.Web.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", app.Conf.Web.Port), router); err != nil {
		log.Fatal(err)
	}
}
