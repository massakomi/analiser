package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
	"time"
)

func (app *application) weeks(w http.ResponseWriter, r *http.Request) {

	var texts []string
	data := weekStatSorted()
	for _, week := range data {
		texts = append(texts, week.category+" "+week.durationFormatted)
	}
	p := struct{ Results []string }{Results: texts}
	app.display("weeks.page.html", w, p)
}

func (app *application) total(w http.ResponseWriter, r *http.Request) {

	data := getData()
	stat := times{}
	for _, dayinfo := range data {
		stat = dayinfo.sumStat(stat)
	}

	p := struct{ Results []string }{Results: viewStatWeb(stat, "Total")}
	app.display("total.page.html", w, p)
}

// Выводит массив статистики по категориям с заголовком
func viewStatWeb(stat times, title string) []string {
	rows := []string{}
	rows = append(rows, title)
	total := time.Duration(0)
	keys := mapKeySortedByValues(stat)
	for _, category := range keys {
		duration := stat[category]
		if category == "" {
			category = "--"
			continue
		} else {
			total += duration
		}
		rows = append(rows, fmtDuration(duration)+" "+category)
	}
	rows = append(rows, "Total: "+fmtDuration(total))
	return rows
}

func (app *application) days(w http.ResponseWriter, r *http.Request) {

	var texts []string
	data := getData()
	for _, dayinfo := range data {
		texts = append(texts, fmt.Sprint(dayinfo))
	}

	p := struct{ Results []string }{Results: texts}
	app.display("days.page.html", w, p)
}

func (app *application) category(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	var texts []string

	rows := statCategory(vars["name"])
	for _, stat := range rows {
		values := []string{}
		for _, period := range stat.periods {
			values = append(values, "["+fmtDuration(period.minutes())+"] "+period.value)
		}
		texts = append(texts, stat.day)
		texts = append(texts, strings.Join(values, "\n"))
		texts = append(texts, strings.Repeat("-", 100))
	}

	p := struct {
		Results  []string
		Category string
	}{
		Results:  texts,
		Category: vars["name"],
	}
	app.display("category.page.html", w, p)
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	app.display("home.page.html", w, nil)
}
