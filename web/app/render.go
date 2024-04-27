package app

import (
	"github.com/rwirdemann/databasedragon/web/templates"
	"html/template"
	"net/http"
)

// MsgSuccess holds a success messages that is shown on the index page.
var MsgSuccess string

// MsgError holds an error message that is shown on the index page.
var MsgError string

type ViewData struct {
	Title   string
	Message string
	Error   string
}

// Render renders tmpl embedded in layout.html using the provided data.
func Render(tmpl string, w http.ResponseWriter, data any) {
	if err := RenderE(tmpl, w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// RenderE works the same as Render except returning the error instead of
// handling it.
func RenderE(tmpl string, w http.ResponseWriter, data any) error {
	t, err := template.ParseFS(templates.Templates, "layout.html", "messages.html", tmpl)
	if err != nil {
		return err
	}
	return t.Execute(w, data)
}

// RenderS renders tmpl embedded in layout.html and inserts title.
func RenderS(tmpl string, w http.ResponseWriter, title string) {
	if err := RenderE(tmpl, w, ViewData{
		Title:   title,
		Message: "",
		Error:   "",
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// RedirectE redirects to url after setting the global msgError to err.
func RedirectE(w http.ResponseWriter, r *http.Request, url string, err error) {
	MsgError = err.Error()
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func ClearMessages() (string, string) {
	e := MsgError
	m := MsgSuccess
	MsgError = ""
	MsgSuccess = ""
	return m, e
}
