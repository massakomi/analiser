package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type dayinfo struct {
	day  string
	text string
}

type period struct {
	prevTime string
	time     string
	value    string
}

type times map[string]time.Duration

type filter struct {
	day string
}

// Сначала разбираем на дату и строки под ней
func getDataFilter(filter filter) []dayinfo {

	const filename = "F:/Google Диск/Задачи.txt"

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	dateTitle := ""
	lines := []string{}
	data := []dayinfo{}
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			if dateTitle != "" {
				data = append(data, dayinfo{dateTitle, strings.Join(lines, " ")})
				if filter.day != "" && filter.day == dateTitle {
					break
				}
				dateTitle = ""
				lines = []string{}
			}
			continue
		}
		// собираем массив строк файла
		if dateTitle != "" {
			lines = append(lines, text)
		}
		// находим dateTitle
		re := regexp.MustCompile(`^\d{1,2}[.,]\d\d`)
		str := re.FindString(text)
		if str != "" {
			dateTitle = strings.Replace(str, ",", ".", 1)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return data
}

func getData() []dayinfo {
	return getDataFilter(filter{})
}

// ---------------------------------------------------------------------------------------------------------------------
// Функции разбора блока текста по дню

// Возвращает массив периодов в указанном дне
func (info dayinfo) getTimeValues() (data []period) {
	re := regexp.MustCompile(`(^| )\d{1,2}-\d\d`)
	matches := re.FindAllStringIndex(info.text, -1)
	data = []period{}
	prevTime := ""
	for index, match := range matches {
		time := strings.TrimSpace(info.text[match[0]:match[1]])
		value := info.getTextForTime(index, match, matches)
		if value != "" {
			data = append(data, period{prevTime, time, value})
		}
		prevTime = time
	}
	return
}

// getTimeValues без пустых категорий
func (info dayinfo) getTimeValuesWithoutEmptyCategory() []period {
	dayValues := info.getTimeValues()
	data := []period{}
	for _, period := range dayValues {
		category := period.category()
		if category == "" {
			continue
		}
		data = append(data, period)
	}
	return data
}

// getTimeValues но только периоды с указанной категорией
func (info dayinfo) getTimeValuesCategory(selectedCategory string) []period {
	dayValues := info.getTimeValues()
	data := []period{}
	for _, period := range dayValues {
		category := period.category()
		if category != selectedCategory {
			continue
		}
		data = append(data, period)
	}
	return data
}

// Получаем текст, относящийся к временному отрезку из текста info.text
func (info dayinfo) getTextForTime(index int, match []int, matches [][]int) string {
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
func (info dayinfo) sumStat(stat times) times {
	dayValues := info.getTimeValues()
	for _, period := range dayValues {
		minutes := period.minutes()
		category := period.category()
		stat[category] += minutes
	}
	return stat
}

// Всего времени отработано за день
func (info dayinfo) total() time.Duration {
	total := time.Duration(0)
	dayValues := info.getTimeValues()
	for _, period := range dayValues {
		minutes := period.minutes()
		if period.category() == "" {
			continue
		}
		total += minutes
	}
	return total
}

// ---------------------------------------------------------------------------------------------------------------------
// Функции конкретного периода времени внутри дня

// Из строки вида HH:mm получаем интервал
func (p period) parseHoursAndMinutes(str string) time.Duration {
	var hours, minutes time.Duration
	_, err := fmt.Sscanf(str, "%d-%d", &hours, &minutes)
	if err != nil {
		log.Fatal(err)
	}
	return time.Hour*hours + time.Minute*minutes
}

// Получаем количество минут между prevTime и time
func (p period) minutes() time.Duration {
	if p.prevTime == "" || p.time == "" {
		return 0
	}
	duration := p.parseHoursAndMinutes(p.prevTime)
	durationTo := p.parseHoursAndMinutes(p.time)
	if duration > durationTo {
		durationTo += time.Hour * 24
	}
	return durationTo - duration
}

// Извлекаем категорию на основе текста периода
func (p period) category() string {
	cat := ""
	categories := []string{"go", "work", "dev", "sql", "read", "python", "php", "par"}
	for _, category := range categories {
		if strings.HasPrefix(p.value, category) {
			cat = category
			break
		}
	}
	return cat
}

// Минуты преобразовать в строку
func (p period) minutesString() string {
	minutes := p.minutes().Minutes()
	return strconv.FormatInt(int64(minutes), 10)
}

// ---------------------------------------------------------------------------------------------------------------------
// Helpers

// Форматирует время
func fmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02d:%02d", h, m)
}

// Сортирует ключи словаря
func mapKeySortedByValues(stat times) []string {
	keys := make([]string, 0, len(stat))
	for key := range stat {
		keys = append(keys, key)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return stat[keys[i]] > stat[keys[j]]
	})
	return keys
}

func getFirstElement[T any](s []T) T {
	return s[0]
}

// ---------------------------------------------------------------------------------------------------------------------
// Анализ категории

type dayinfoEx struct {
	day     string
	periods []period
}

// Возвращаем массив дней и периодов (текстов) только по конкретной категории
// По сути это фильтрованный getData состоящий не из dayinfo а из dayinfoEx
func statCategory(selectedCategory string) (rows []dayinfoEx) {
	rows = []dayinfoEx{}
	data := getData()
	for _, dayinfo := range data {
		values := dayinfo.getTimeValuesCategory(selectedCategory)
		if len(values) > 0 {
			rows = append(rows, dayinfoEx{dayinfo.day, values})
		}
	}
	return
}

// ---------------------------------------------------------------------------------------------------------------------
// Статистика за конкретный день
func getDayInfoByDate(date string) dayinfo {
	if date == "" {
		date = time.Now().Add(-6 * time.Hour).Format("02.01")
	}
	data := getDataFilter(filter{day: date})
	for _, dayinfo := range data {
		if dayinfo.day == date {
			return dayinfo
		}
	}
	return dayinfo{}
}

// ---------------------------------------------------------------------------------------------------------------------
func main() {
	osArgs := os.Args[1:]
	if len(osArgs) > 0 {
		consoleProcess()
	} else {
		webProcess()
	}
}
