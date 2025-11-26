package web

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"runtime/debug"
)

// Routing (using gorilla/mux)
// простейший пример https://gowebexamples.com/routes-using-gorilla-mux/
// документация https://github.com/gorilla/mux?tab=readme-ov-file

type application struct {
}

func Process() {

	app := &application{}

	r := mux.NewRouter()
	r.HandleFunc("/", app.home)
	r.HandleFunc("/weeks", app.weeks)
	r.HandleFunc("/days", app.days)
	r.HandleFunc("/total", app.total)
	r.HandleFunc("/category/{name}", app.category)

	handler := http.StripPrefix("/static/", http.FileServer(http.Dir("./static")))
	r.PathPrefix("/static/").Handler(handler)

	log.Println("Запуск веб-сервера на http://127.0.0.1:4000")
	err := http.ListenAndServe(":4000", r)

	log.Fatal(err)
}

func (app *application) display(tpl string, w http.ResponseWriter, data any) {
	files := []string{
		"./templates/" + tpl,
		//"./templates/home.page.html",
		"./templates/base.layout.html",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	log.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}
