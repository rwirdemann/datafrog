package app

import "github.com/gorilla/mux"

func RegisterHandler(router *mux.Router) {
	// home
	router.HandleFunc("/", IndexHandler)

	// show new form
	router.HandleFunc("/new", NewHandler)

	router.HandleFunc("/create", CreateHandler)

	// delete test
	router.HandleFunc("/delete", DeleteHandler)

	// start recording
	router.HandleFunc("/record", StartRecording)

	// stop recording
	router.HandleFunc("/stoprecording", StopRecording)

	// start verifx
	router.HandleFunc("/run", StartHandler)

	// stop verify
	router.HandleFunc("/stop", StopHandler)

	// show test
	router.HandleFunc("/show", ShowHandler)

}
