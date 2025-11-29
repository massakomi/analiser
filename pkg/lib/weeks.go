package lib

import (
	"maps"
	"slices"
	"strconv"
	"time"
)

// Анализ по неделям
type weekStat struct {
	Category          string
	Duration          time.Duration
	DurationFormatted string
}

// Уже подготвленный для конечного вывода массив
func WeekStatSorted() []weekStat {
	stat := groupWeekStat()
	keys := slices.Sorted(maps.Keys(stat))
	slices.Reverse(keys)
	var data []weekStat
	for _, category := range keys {
		duration := stat[category]
		if category == "" {
			category = "--"
			continue
		}
		data = append(data, weekStat{category, duration, FmtDuration(duration)})
	}
	return data
}

// Простая несортированная структура с пустыми категориями
func groupWeekStat() map[string]time.Duration {
	data := GetData()
	stat := map[string]time.Duration{}
	prevDay := ""
	prevWeek := ""
	slices.Reverse(data)
	for _, dayinfo := range data {
		week := dayinfo.WeekNum()
		weekStr := strconv.FormatInt(int64(week), 10)
		if prevDay == "" || weekStr != prevWeek {
			prevDay = dayinfo.Day
		}
		prevWeek = weekStr
		stat[weekStr+"/"+prevDay] += dayinfo.Total()
	}
	return stat
}

// Текущий номер недели в году у этого дня
var currentYear = strconv.FormatInt(int64(time.Now().Year()), 10)

func (info Dayinfo) WeekNum() (w int) {
	input := currentYear + "." + info.Day
	const layout = "2006.02.01"
	t, err := time.Parse(layout, input)
	if err != nil {
		panic(err)
	}
	_, w = t.ISOWeek()
	return
}
