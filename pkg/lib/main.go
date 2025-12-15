package lib

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

type Dayinfo struct {
	Day  string
	text string
}

type Period struct {
	prevTime string
	time     string
	Value    string
}

type Times map[string]time.Duration

func taskLines(yield func(string) bool) {
	const filename = "F:/Google Диск/Задачи.txt"

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		yield(scanner.Text())
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
}

// Сначала разбираем на дату и строки под ней
func GetData() []Dayinfo {
	dateTitle := ""
	lines := []string{}
	data := []Dayinfo{}
	for text := range taskLines {
		if text == "" {
			if dateTitle != "" {
				data = append(data, Dayinfo{dateTitle, strings.Join(lines, " ")})
				/*if filter.day != "" && filter.day == dateTitle {
					break
				}*/
				dateTitle = ""
				lines = []string{}
			}
			continue
		}
		// собираем массив строк файла
		if dateTitle != "" {
			lines = append(lines, text)
		}
		findDateTitle(text, &dateTitle)
	}
	return data
}

// находим dateTitle
func findDateTitle(text string, dateTitle *string) {
	re := regexp.MustCompile(`^\d{1,2}[.,]\d\d`)
	str := re.FindString(text)
	if str == "" {
		return
	}
	str = strings.TrimSpace(str)
	if len(str) != 5 {
		str = "0" + str
	}
	if str != "" {
		*dateTitle = strings.Replace(str, ",", ".", 1)
	}
}

func Last7days() []string {
	data := GetData()
	var lastDays []string
	for _, item := range data {
		lastDays = append(lastDays, item.Day)
		if len(lastDays) == 7 {
			break
		}
	}
	return lastDays
}

// ---------------------------------------------------------------------------------------------------------------------
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

// ---------------------------------------------------------------------------------------------------------------------
// Функции конкретного периода времени внутри дня

// Из строки вида HH:mm получаем интервал
func (p Period) parseHoursAndMinutes(str string) time.Duration {
	var hours, minutes time.Duration
	_, err := fmt.Sscanf(str, "%d-%d", &hours, &minutes)
	if err != nil {
		log.Fatal(err)
	}
	return time.Hour*hours + time.Minute*minutes
}

// Получаем количество минут между prevTime и time
func (p Period) Minutes() time.Duration {
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
func (p Period) Category() string {
	cat := ""
	for _, category := range Categories {
		if strings.HasPrefix(p.Value, category) {
			cat = category
			break
		}
	}
	return cat
}

var Categories = []string{"par", "go", "python", "dev", "work", "sql", "read", "php"}

// Минуты преобразовать в строку
func (p Period) MinutesString() string {
	minutes := p.Minutes().Minutes()
	return strconv.FormatInt(int64(minutes), 10)
}

// ---------------------------------------------------------------------------------------------------------------------
// Helpers

// Форматирует время
func FmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02d:%02d", h, m)
}

// Сортирует ключи словаря
func MapKeySortedByValues(stat Times) []string {
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

type DayinfoEx struct {
	Day     string
	Periods []Period
}

// Возвращаем массив дней и периодов (текстов) только по конкретной категории
// По сути это фильтрованный getData состоящий не из Dayinfo а из DayinfoEx
func StatCategory(selectedCategory string) (rows []DayinfoEx) {
	rows = []DayinfoEx{}
	data := GetData()
	for _, Dayinfo := range data {
		values := Dayinfo.getTimeValuesCategory(selectedCategory)
		if len(values) > 0 {
			rows = append(rows, DayinfoEx{Dayinfo.Day, values})
		}
	}
	return
}

// ---------------------------------------------------------------------------------------------------------------------
// Статистика за конкретный день
func GetDayinfoByDate(date string) Dayinfo {
	if date == "" {
		date = time.Now().Add(-6 * time.Hour).Format("02.01")
	}
	data := GetData()
	for _, Dayinfo := range data {
		if Dayinfo.Day == date {
			return Dayinfo
		}
	}
	return Dayinfo{}
}
