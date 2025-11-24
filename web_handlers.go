package main

import "net/http"

func (app *application) weeks(w http.ResponseWriter, r *http.Request) {

	var texts []string
	data := weekStatSorted()
	for _, week := range data {
		texts = append(texts, week.category+" "+week.durationFormatted)
	}
	p := struct{ Results []string }{Results: texts}
	app.display("weeks.page.html", w, p)
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	app.display("home.page.html", w, nil)
}
