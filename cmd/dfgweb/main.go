package main

import (
	"embed"
	"github.com/rwirdemann/datafrog/pkg/df"
	"github.com/rwirdemann/datafrog/pkg/web"
	"github.com/rwirdemann/simpleweb"
	"log"
)

// Expects all HTML templates in datafrog/cmd/dfgweb/templates
//
//go:embed all:templates
var templates embed.FS

func main() {
	config, err := df.NewDefaultConfig()
	if err != nil {
		log.Fatal(err)
	}
	simpleweb.Init(templates, []string{"templates/layout.html"}, config.Web.Port)
	web.RegisterHandler(config)
	simpleweb.ShowRoutes()
	log.Printf("Connecting backend: http://localhost:/%d", config.Api.Port)
	simpleweb.Run()
}
