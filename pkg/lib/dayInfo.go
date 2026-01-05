package lib

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Функции разбора блока текста по дню

// Возвращает массив периодов в указанном дне
func (info Dayinfo) getTimeValues() (data []Period) {
	re := regexp.MustCompile(`(^| )\d{1,2}-\d\d`)
	matches := re.FindAllStringIndex(info.text, -1)
	data = []Period{}
	prevTime := ""
	for index, match := range matches {
		time := strings.TrimSpace(info.text[match[0]:match[1]])
		value := info.getTextForTime(index, match, matches)
		if value != "" {
			data = append(data, Period{prevTime, time, value})
		}
		prevTime = time
	}
	return
}

// getTimeValues без пустых категорий
func (info Dayinfo) GetTimeValuesWithoutEmptyCategory() []Period {
	dayValues := info.getTimeValues()
	data := []Period{}
	for _, Period := range dayValues {
		category := Period.Category()
		if category == "" {
			continue
		}
		data = append(data, Period)
	}
	return data
}

// getTimeValues но только периоды с указанной категорией
func (info Dayinfo) getTimeValuesCategory(selectedCategory string) []Period {
	dayValues := info.getTimeValues()
	data := []Period{}
	for _, Period := range dayValues {
		category := Period.Category()
		if category != selectedCategory {
			continue
		}
		data = append(data, Period)
	}
	return data
}

// Получаем текст, относящийся к временному отрезку из текста info.text
func (info Dayinfo) getTextForTime(index int, match []int, matches [][]int) string {
	value := ""
	nextIndex := index + 1
	if len(matches)-1 >= nextIndex {
		nextMatch := matches[nextIndex]
		value = info.text[match[1]:nextMatch[0]]
	} else {
		value = info.text[match[1]:]
	}
	if len(value) > 0 {
		value = value[1:]
	}
	return value
}

// Суммируем статистику за день (категория - время)
func (info Dayinfo) SumStat(stat Times) Times {
	dayValues := info.getTimeValues()
	for _, Period := range dayValues {
		minutes := Period.Minutes()
		category := Period.Category()
		stat[category] += minutes
	}
	return stat
}

// Всего времени отработано за день
func (info Dayinfo) Total() time.Duration {
	total := time.Duration(0)
	dayValues := info.getTimeValues()
	for _, Period := range dayValues {
		minutes := Period.Minutes()
		if Period.Category() == "" {
			continue
		}
		total += minutes
	}
	return total
}

// Stringer() интерфейс (как магический метод), можно структуру вывести просто строкой  fmt.Println(dayinfo)
func (info Dayinfo) String() string {
	total := info.Total()
	graph := strings.Repeat("-", int(total.Minutes())/5)
	return fmt.Sprintf("%v %v %v", info.Day, FmtDuration(total), graph)
}
