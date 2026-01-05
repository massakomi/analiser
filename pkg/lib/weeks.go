package lib

import (
	"maps"
	"slices"
	"strconv"
	"time"
)

// Анализ по неделям
type WeekStat struct {
	Category          string
	Duration          time.Duration
	DurationFormatted string
}

// Уже подготвленный для конечного вывода массив
func WeekStatSorted() []WeekStat {
	stat := groupWeekStat()
	keys := slices.Sorted(maps.Keys(stat))
	slices.Reverse(keys)
	var data []WeekStat
	for _, category := range keys {
		duration := stat[category]
		if category == "" {
			category = "--"
			continue
		}
		data = append(data, WeekStat{category, duration, FmtDuration(duration)})
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
		year := strconv.FormatInt(int64(dayinfo.Year), 10)
		if prevDay == "" || weekStr != prevWeek {
			prevDay = dayinfo.Day
		}
		prevWeek = weekStr
		stat[year+" / "+weekStr+" / "+prevDay] += dayinfo.Total()
	}
	return stat
}

// Текущий номер недели в году у этого дня
var currentYearStr = strconv.FormatInt(int64(time.Now().Year()), 10)

func (info Dayinfo) WeekNum() (w int) {
	input := currentYearStr + "." + info.Day
	const layout = "2006.02.01"
	t, err := time.Parse(layout, input)
	if err != nil {
		panic(err)
	}
	_, w = t.ISOWeek()
	return
}

// Преобразование массива WeekStat в простой словарь с процентными данными для вывода в шаблон
func WeekStatTexts(data []WeekStat) []map[string]string {
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
			"text":    week.Category,
			"time":    week.DurationFormatted,
			"percent": strconv.FormatFloat(percent, 'f', 1, 32),
		})
	}
	return texts
}
