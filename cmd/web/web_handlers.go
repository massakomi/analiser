package web

import (
	"analiser/pkg/lib"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	// по сути тут то же самое, что и day, только за 7 дней, можно сделать одну функцию
	fullRows := []string{}
	data := lib.GetData()[0:7]
	for _, info := range data {
		fullRows = append(fullRows, dayRowsByDayInfo(info)...)
		stat := lib.Times{}
		stat = info.SumStat(stat)
		fullRows = append(fullRows, strings.Join(viewStatWeb(stat), " "))
	}

	p := map[string]any{
		"data": fullRows,
	}
	app.display("home.page.html", w, p)
}

func (app *application) day(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	info := lib.GetDayinfoByDate(vars["date"])
	rows := dayRowsByDayInfo(info)

	stat := lib.Times{}
	stat = info.SumStat(stat)

	p := map[string]any{
		"data": append(rows, viewStatWeb(stat)...),
		"date": vars["date"],
	}
	app.display("day.page.html", w, p)
}

func (app *application) weeks(w http.ResponseWriter, r *http.Request) {
	data := lib.WeekStatSorted()

	maxDuration := time.Duration(0)
	for _, week := range data {
		if week.Duration > maxDuration {
			maxDuration = week.Duration
		}
	}

	texts := []map[string]string{}
	for _, week := range data {
		percent := 100 * (week.Duration.Seconds() / maxDuration.Seconds())
		texts = append(texts, map[string]string{
			"text":    week.Category + " " + week.DurationFormatted,
			"percent": strconv.FormatFloat(percent, 'f', 1, 32),
		})
	}

	p := map[string]any{
		"data": texts,
	}
	app.display("weeks.page.html", w, p)
}

func (app *application) total(w http.ResponseWriter, r *http.Request) {

	data := lib.GetData()
	stat := lib.Times{}
	for _, dayinfo := range data {
		stat = dayinfo.SumStat(stat)
	}
	p := map[string]any{
		"data": viewStatWeb(stat),
	}
	app.display("total.page.html", w, p)
}

func (app *application) days(w http.ResponseWriter, r *http.Request) {

	var texts []string
	data := lib.GetData()
	for _, dayinfo := range data {
		texts = append(texts, fmt.Sprint(dayinfo))
	}
	p := map[string]any{
		"data": texts,
	}
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

	p := map[string]any{
		"data":     texts,
		"category": vars["name"],
	}
	app.display("category.page.html", w, p)
}
