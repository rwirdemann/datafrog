package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rwirdemann/datafrog/pkg/api"
	"github.com/rwirdemann/datafrog/pkg/df"
	"github.com/rwirdemann/datafrog/pkg/file"
	"log"
	"net/http"
)

func main() {
	router := mux.NewRouter()
	verificationDoneChannels := make(api.ChannelMap)
	verificationStoppedChannels := make(api.ChannelMap)
	recordingDoneChannels := make(api.ChannelMap)
	recordingStoppedChannels := make(api.ChannelMap)

	config, err := df.NewDefaultConfig()
	if err != nil {
		log.Fatal(err)
	}
	testRepository := file.TestRepository{}
	api.RegisterHandler(config, router,
		verificationDoneChannels, verificationStoppedChannels,
		recordingDoneChannels, recordingStoppedChannels, testRepository)
	err = router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, _ := route.GetPathTemplate()
		met, _ := route.GetMethods()
		log.Println(tpl, met)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening on :%d...", config.Api.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", config.Api.Port), router); err != nil {
		log.Fatal(err)
	}
}
