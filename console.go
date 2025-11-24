package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

func consoleProcess() {
	var option = flag.String("o", "", "sub option for script")
	var date = flag.String("date", "", "date for analizer")
	var category = flag.String("c", "", "category for analizer")
	flag.Parse()
	analize(*date, *option, *category)
}

// ---------------------------------------------------------------------------------------------------------------------
// Точка входа для анализа
func analize(date, option, category string) {
	if option == "total" {
		printTotalConsole()
	} else if option == "days" {
		printAllDaysGraphConsole()
	} else if option == "weeks" {
		printWeeksConsole()
	} else if category != "" {
		printCategoryConsole(category)
	} else {
		printDateStatConsole(date)
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Статистика по всем дням в виде графика
func printAllDaysGraphConsole() {
	data := getData()
	for _, dayinfo := range data {
		fmt.Println(dayinfo)
	}
}

// Stringer() интерфейс (как магический метод), можно структуру вывести просто строкой  fmt.Println(dayinfo)
func (info dayinfo) String() string {
	total := info.total()
	graph := strings.Repeat("-", int(total.Minutes())/5)
	return fmt.Sprintf("%v %v %v", info.day, fmtDuration(total), graph)
}

// ---------------------------------------------------------------------------------------------------------------------
// Анализ категории
func printCategoryConsole(selectedCategory string) {
	rows := statCategory(selectedCategory)
	for _, stat := range rows {
		values := []string{}
		for _, period := range stat.periods {
			values = append(values, "["+fmtDuration(period.minutes())+"] "+period.value)
		}
		fmt.Println(stat.day)
		fmt.Println(strings.Join(values, "\n"))
		fmt.Println(strings.Repeat("-", 100))
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// По неделям
func printWeeksConsole() {
	data := weekStatSorted()
	for _, week := range data {
		fmt.Println(week.category + " " + week.durationFormatted)
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Суммарное количество времени за все дни
func printTotalConsole() {
	data := getData()
	stat := times{}
	for _, dayinfo := range data {
		stat = dayinfo.sumStat(stat)
	}
	viewStatConsole(stat, "Total")
}

// ---------------------------------------------------------------------------------------------------------------------
// Анализ за конкретный день
func printDateStatConsole(date string) {
	info := getDayInfoByDate(date)
	data := info.getTimeValuesWithoutEmptyCategory()

	fmt.Println(date)
	fmt.Println(strings.Repeat("=", 60))
	for _, period := range data {
		fmt.Println("["+period.minutesString()+"]", period.value)
	}
	stat := times{}
	stat = info.sumStat(stat)
	viewStatConsole(stat, info.day)
}

// ---------------------------------------------------------------------------------------------------------------------
// Выводит массив статистики по категориям с заголовком
func viewStatConsole(stat times, title string) {
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println(title)
	total := time.Duration(0)
	keys := mapKeySortedByValues(stat)
	for _, category := range keys {
		duration := stat[category]
		if category == "" {
			category = "--"
			continue
		} else {
			total += duration
		}
		fmt.Printf(warn()+" "+fmtDuration(duration)+" ", category)
	}
	fmt.Printf("\n"+warn()+" ", "Total:")
	fmt.Println(fmtDuration(total))
}

func warn() string {
	if hasArg("-o=nocolor") {
		return "%v"
	}
	var reset = "\033[0m"
	var red = "\033[31m"
	return red + "%v" + reset
}

// Проверяет указан ли этот аргумент в командной строке
func hasArg(arg string) bool {
	return slices.Contains(os.Args, arg)
}
