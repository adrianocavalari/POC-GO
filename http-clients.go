package main

import (
	"flag"
	"fmt"
	htemplate "html/template"
	"net/http"
	"text/template"
	"time"

	"github.com/aarol/reload"
)

var tmpl *template.Template

var handler http.Handler = http.DefaultServeMux

func main() {
	isDevelopment := flag.Bool("dev", true, "Development mode")
	templateCache := parseTemplates()

	// handler can be anything that implements http.Handler,
	// like chi.Router, echo.Echo or gin.Engine
	var handler http.Handler = http.DefaultServeMux

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// serve a template file with dynamic data
		data := map[string]any{
			"Timestamp": time.Now().Format("Monday, 02-Jan-06 15:04:05 MST"),
		}
		err := templateCache.ExecuteTemplate(w, "index.html", data)
		if err != nil {
			fmt.Println(err)
		}
	})

	if *isDevelopment {
		// Call `New()` with a list of directories to recursively watch
		reload := reload.New("ui/")

		reload.OnReload = func() {
			templateCache = parseTemplates()
		}

		handler = reload.Handle(handler)
	} else {
		fmt.Println("Running in production mode")
	}

	tmpl = template.Must(template.ParseFiles("index.html"))
	http.HandleFunc("/", foo)
	http.ListenAndServe(":3000", nil)
}

func foo(reswt http.ResponseWriter, req *http.Request) {
	tmpl.ExecuteTemplate(reswt, "index.html", nil)
}

func parseTemplates() *htemplate.Template {
	return htemplate.Must(htemplate.ParseGlob("ui/*.html"))
}

// output html
func OutputHTML(w http.ResponseWriter, filename string, data interface{}) {
	t, err := template.ParseFiles(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
