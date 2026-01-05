package web

import (
	"analiser/pkg/lib"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (app *application) table(w http.ResponseWriter, r *http.Request) {
	var dates []map[string]string
	table := map[string]map[string]map[string]any{}
	all := lib.GetData()
	for _, dayinfo := range all {
		if table[dayinfo.Day] == nil {
			table[dayinfo.Day] = map[string]map[string]any{}
		}
		data := dayinfo.GetTimeValuesWithoutEmptyCategory()

		// Собираем все даты, время, тексты
		dates = append(dates, map[string]string{
			"date":  dayinfo.Day,
			"month": dayinfo.Month,
			"year":  strconv.Itoa(dayinfo.Year),
			"day":   strings.TrimLeft(dayinfo.Day[0:2], "0"),
		})
		sums, texts := collectStatsForTable(data)

		// Итоговый общий массив со всеми данными
		for category, minutes := range sums {
			table[dayinfo.Day][category] = map[string]any{
				"minutes": minutes.Minutes(),
				"text":    strings.Join(texts[category], " / "),
			}
		}
	}

	p := map[string]any{
		"data":    table,
		"monthes": getMonthes(dates),
		"dates":   dates,
	}
	app.display("table.page.html", w, p)
}

func getMonthes(dates []map[string]string) []map[string]int {
	var monthes []map[string]int
	var count int
	var countAll = len(dates)
	for key, item := range dates {
		count++
		if key+1 != countAll {
			nextMonth := dates[key+1]["month"]
			if item["month"] != nextMonth {
				monthName := getMonthName(item["month"])
				monthes = append(monthes, map[string]int{
					monthName: count,
				})
				count = 0
			}
		} else {
			monthName := getMonthName(item["month"])
			monthes = append(monthes, map[string]int{
				monthName: count,
			})
		}
	}
	return monthes
}

func getMonthName(index string) string {
	monthes := []string{"Январь", "Февраль", "Март", "Апрель", "Май", "Июнь", "Июль", "Август", "Сентябрь", "Октябрь", "Ноябрь", "Декабрь"}
	i, _ := strconv.Atoi(index)
	return monthes[i-1]
}

func collectStatsForTable(data []lib.Period) (map[string]time.Duration, map[string][]string) {
	sums := map[string]time.Duration{}
	texts := map[string][]string{}
	for _, period := range data {

		// Собираем время по каждой категории по текущему дню
		sums[period.Category()] += period.Minutes()

		// Собираем тексты для подсказок
		line, _ := strings.CutPrefix(period.Value, period.Category()+" ")
		maxLength := 30
		runes := []rune(line)
		if len(runes) > maxLength {
			line = string(runes[:maxLength]) + "..."
		}
		texts[period.Category()] = append(texts[period.Category()], line)
	}
	return sums, texts
}
