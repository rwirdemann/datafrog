package cmd

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rwirdemann/databasedragon/http/api"
	"github.com/rwirdemann/databasedragon/web/app"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

func init() {
	rootCmd.AddCommand(apiCmd)
}

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the backend API",
	Run: func(cmd *cobra.Command, args []string) {
		router := mux.NewRouter()
		api.RegisterHandler(router)
		err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			tpl, _ := route.GetPathTemplate()
			met, _ := route.GetMethods()
			log.Println(tpl, met)
			return nil
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Listening on :%d...", app.Conf.Api.Port)
		if err := http.ListenAndServe(fmt.Sprintf(":%d", app.Conf.Api.Port), router); err != nil {
			log.Fatal(err)
		}
	},
}
