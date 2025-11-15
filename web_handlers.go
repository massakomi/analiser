package main

import "net/http"

type Page struct {
	Content string
	Results []string
}

func (app *application) weeks(w http.ResponseWriter, r *http.Request) {

	var texts []string
	data := weekStatSorted()
	for _, week := range data {
		texts = append(texts, week.category+" "+week.durationFormatted)
	}

	p := &Page{
		Results: texts,
	}

	app.display("weeks.page.html", w, p)
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	app.display("home.page.html", w, nil)
}
