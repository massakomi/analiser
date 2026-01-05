package web

import (
	"analiser/pkg/lib"
	"time"
)

// Выводит массив статистики по категориям с заголовком
func viewStatWeb(stat lib.Times) []string {
	rows := []string{}
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

func dayRowsByDayInfo(info lib.Dayinfo) []string {
	data := info.GetTimeValuesWithoutEmptyCategory()
	rows := []string{}
	//rows = append(rows, strings.Repeat("=", 60))
	rows = append(rows, info.Day)
	//rows = append(rows, strings.Repeat("=", 60))
	for _, period := range data {
		rows = append(rows, "["+period.MinutesString()+"] "+period.Value)
	}
	return rows
}
