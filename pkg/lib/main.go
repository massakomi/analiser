package lib

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Dayinfo struct {
	Day   string
	text  string
	Month string
	Year  int
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
				month, year := getMonthYear(data, dateTitle)
				data = append(data, Dayinfo{dateTitle, strings.Join(lines, " "), month, year})
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

var currentYear = time.Now().Year()

// Извлекает текущий месяц и год из строки данных вида "ДД.ММ" "30.12"
func getMonthYear(data []Dayinfo, dateTitle string) (string, int) {
	month := ""
	year := currentYear
	month = dateTitle[strings.Index(dateTitle, ".")+1:]
	if len(data) > 0 {
		nextDateTitle := data[len(data)-1].Day
		year = data[len(data)-1].Year
		nextMonth := nextDateTitle[strings.Index(nextDateTitle, ".")+1:]
		if nextMonth < month {
			year--
		}
	}
	return month, year
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
// Функции конкретного периода времени внутри дня

type Period struct {
	prevTime string
	time     string
	Value    string
}

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
	categories := Categories()
	spaceIndex := strings.Index(p.Value, " ")
	if spaceIndex > 0 {
		cat := p.Value[:spaceIndex]
		if slices.Contains(categories, cat) {
			return cat
		}
	}
	return ""
}

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
func MapKeySortedByValuesAssoc(stat map[string]int) []string {
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

func AppendToFile(filename string, content string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err = f.WriteString(content); err != nil {
		panic(err)
	}
}

func WriteFile(filename string, content []byte) {
	err := os.WriteFile(filename, content, 0666)
	if err != nil {
		log.Fatal(err)
	}
}

func sliceUnique(strings []string) []string {
	slices.Sort(strings)
	strings = slices.Compact(strings)
	return strings
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

var categories = []string{}

// Доступные категории
func Categories() []string {
	if len(categories) == 0 {
		filename := "data/categories.json"
		jsonData, err := os.ReadFile(filename)
		if err != nil {
			log.Fatalf(`File %v not found`, filename)
		}
		err = json.Unmarshal(jsonData, &categories)
		if err != nil {
			log.Fatalf(`Error parsing json categories: %v`, jsonData)
		}
	}
	return categories
}

func CategoryAdd(category string) {
	categories := Categories()
	if !slices.Contains(categories, category) {
		categories = append(categories, category)
		CategoriesSave(categories)
	}
}

func CategoryDelete(name string) {
	categories := Categories()
	if slices.Contains(categories, name) {
		current := slices.Index(categories, name)
		categories = slices.Delete(categories, current, current+1)
		CategoriesSave(categories)
	}
}

func CategoryEdit(name string, currentName string) {
	categories := Categories()
	if slices.Contains(categories, currentName) {
		current := slices.Index(categories, currentName)
		categories = slices.Delete(categories, current, current+1)
		categories = append(categories, name)
		CategoriesSave(categories)
	}
}

func CategoriesSave(data []string) {
	filename := "data/categories.json"
	content, _ := json.MarshalIndent(data, "", "  ")
	WriteFile(filename, content)
}

// ---------------------------------------------------------------------------------------------------------------------
// Статистика за конкретный день
func GetDayinfoByDate(date string) Dayinfo {
	if date == "" {
		date = time.Now().Add(-6 * time.Hour).Format("02.01")
	} else if len(date) == 4 {
		date = "0" + date
	}
	data := GetData()
	for _, Dayinfo := range data {
		if Dayinfo.Day == date {
			return Dayinfo
		}
	}
	return Dayinfo{}
}
