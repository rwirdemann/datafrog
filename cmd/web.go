package cmd

import (
	"embed"
	"github.com/rwirdemann/datafrog/web/app"
	"github.com/rwirdemann/simpleweb"
	"github.com/spf13/cobra"
	"log"
)

// Expects all HTML templates in datafrog/cmd/templates
//
//go:embed all:templates
var templates embed.FS

func init() {
	rootCmd.AddCommand(webCmd)
	simpleweb.Init(templates, []string{"templates/layout.html"}, app.Conf.Web.Port)
}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Starts local web application",
	Run: func(cmd *cobra.Command, args []string) {
		app.RegisterHandler()
		simpleweb.ShowRoutes()
		log.Printf("Connecting backend: http://localhost:/%d", app.Conf.Api.Port)
		simpleweb.Run()
	},
}
