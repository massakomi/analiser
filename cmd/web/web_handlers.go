package web

import (
	"analiser/pkg/lib"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	// по сути тут то же самое, что и day, только за 7 дней, можно сделать одну функцию
	var fullRows []map[string]any
	data := lib.GetData()[0:7]
	for _, info := range data {
		rows := dayRowsByDayInfo(info)
		day := rows[0:1][0]
		rows = rows[1:]
		//fullRows = append(fullRows, dayRowsByDayInfo(info)...)
		stat := lib.Times{}
		stat = info.SumStat(stat)
		totals := strings.Join(viewStatWeb(stat), " ")

		fullRows = append(fullRows, map[string]any{
			"day":    day,
			"rows":   rows,
			"totals": totals,
		})
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
	rows = rows[1:]

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
	p := map[string]any{
		"data": lib.WeekStatTexts(data),
	}
	app.display("weeks.page.html", w, p)
}

func (app *application) days(w http.ResponseWriter, r *http.Request) {
	var data []lib.WeekStat
	for _, dayinfo := range lib.GetData() {
		duration := dayinfo.Total()
		data = append(data, lib.WeekStat{dayinfo.Day, duration, lib.FmtDuration(duration)})
	}
	p := map[string]any{
		"data": lib.WeekStatTexts(data),
	}
	app.display("days.page.html", w, p)
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

func (app *application) manage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("action") == "category-add" {
		lib.CategoryAdd(r.URL.Query().Get("name"))
	}
	if r.URL.Query().Get("action") == "category-delete" {
		lib.CategoryDelete(r.URL.Query().Get("name"))
	}
	if r.URL.Query().Get("action") == "category-edit" {
		lib.CategoryEdit(r.URL.Query().Get("name"), r.URL.Query().Get("currentName"))
	}
	p := map[string]any{}
	app.display("manage.page.html", w, p)
}
