package web

import (
	"analiser/pkg/lib"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
	"time"
)

func (app *application) weeks(w http.ResponseWriter, r *http.Request) {

	var texts []string
	data := lib.WeekStatSorted()
	for _, week := range data {
		texts = append(texts, week.Category+" "+week.DurationFormatted)
	}
	p := struct{ Results []string }{Results: texts}
	app.display("weeks.page.html", w, p)
}

func (app *application) total(w http.ResponseWriter, r *http.Request) {

	data := lib.GetData()
	stat := lib.Times{}
	for _, dayinfo := range data {
		stat = dayinfo.SumStat(stat)
	}

	p := struct{ Results []string }{Results: viewStatWeb(stat, "Total")}
	app.display("total.page.html", w, p)
}

// Выводит массив статистики по категориям с заголовком
func viewStatWeb(stat lib.Times, title string) []string {
	rows := []string{}
	rows = append(rows, title)
	total := time.Duration(0)
	keys := lib.MapKeySortedByValues(stat)
	for _, category := range keys {
		duration := stat[category]
		if category == "" {
			category = "--"
			continue
		} else {
			total += duration
		}
		rows = append(rows, lib.FmtDuration(duration)+" "+category)
	}
	rows = append(rows, "Total: "+lib.FmtDuration(total))
	return rows
}

func (app *application) days(w http.ResponseWriter, r *http.Request) {

	var texts []string
	data := lib.GetData()
	for _, dayinfo := range data {
		texts = append(texts, fmt.Sprint(dayinfo))
	}

	p := struct{ Results []string }{Results: texts}
	app.display("days.page.html", w, p)
}

func (app *application) category(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	var texts []string

	rows := lib.StatCategory(vars["name"])
	for _, stat := range rows {
		values := []string{}
		for _, period := range stat.Periods {
			values = append(values, "["+lib.FmtDuration(period.Minutes())+"] "+period.Value)
		}
		texts = append(texts, stat.Day)
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
