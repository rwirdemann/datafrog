package main

import (
	"embed"
	"github.com/rwirdemann/datafrog/pkg/api/web"
	"github.com/rwirdemann/simpleweb"
	"log"
)

// Expects all HTML templates in datafrog/cmd/dfgweb/templates
//
//go:embed all:templates
var templates embed.FS

func init() {
	simpleweb.Init(templates, []string{"templates/layout.html"}, web.Conf.Web.Port)
}

func main() {
	web.RegisterHandler()
	simpleweb.ShowRoutes()
	log.Printf("Connecting backend: http://localhost:/%d", web.Conf.Api.Port)
	simpleweb.Run()
}
