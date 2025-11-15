package main

import (
	"maps"
	"slices"
	"strconv"
	"time"
)

// Анализ по неделям
type weekStat struct {
	category          string
	duration          time.Duration
	durationFormatted string
}

// Уже подготвленный для конечного вывода массив
func weekStatSorted() []weekStat {
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
		data = append(data, weekStat{category, duration, fmtDuration(duration)})
	}
	return data
}

// Простая несортированная структура с пустыми категориями
func groupWeekStat() map[string]time.Duration {
	data := getData()
	stat := map[string]time.Duration{}
	prevDay := ""
	prevWeek := ""
	slices.Reverse(data)
	for _, dayinfo := range data {
		week := dayinfo.weekNum()
		weekStr := strconv.FormatInt(int64(week), 10)
		if prevDay == "" || weekStr != prevWeek {
			prevDay = dayinfo.day
		}
		prevWeek = weekStr
		stat[weekStr+"/"+prevDay] += dayinfo.total()
	}
	return stat
}

// Текущий номер недели в году у этого дня
var currentYear = strconv.FormatInt(int64(time.Now().Year()), 10)

func (info dayinfo) weekNum() int {
	input := currentYear + "." + info.day
	layout := "2006.02.01"
	t, err := time.Parse(layout, input)
	if err != nil {
		panic(err)
	}
	_, w := t.ISOWeek()
	return w
}
