package console

import (
	"analiser/pkg/lib"
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

func Process() {
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
	data := lib.GetData()
	for _, dayinfo := range data {
		fmt.Println(dayinfo)
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Анализ категории
func printCategoryConsole(selectedCategory string) {
	rows := lib.StatCategory(selectedCategory)
	for _, stat := range rows {
		values := []string{}
		for _, period := range stat.Periods {
			values = append(values, "["+lib.FmtDuration(period.Minutes())+"] "+period.Value)
		}
		fmt.Println(stat.Day)
		fmt.Println(strings.Join(values, "\n"))
		fmt.Println(strings.Repeat("-", 100))
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// По неделям
func printWeeksConsole() {
	data := lib.WeekStatSorted()
	for _, week := range data {
		fmt.Println(week.Category + " " + week.DurationFormatted)
	}
}

// ---------------------------------------------------------------------------------------------------------------------
// Суммарное количество времени за все дни
func printTotalConsole() {
	data := lib.GetData()
	stat := lib.Times{}
	for _, dayinfo := range data {
		stat = dayinfo.SumStat(stat)
	}
	viewStatConsole(stat, "Total")
}

// ---------------------------------------------------------------------------------------------------------------------
// Анализ за конкретный день
func printDateStatConsole(date string) {
	info := lib.GetDayinfoByDate(date)
	data := info.GetTimeValuesWithoutEmptyCategory()

	fmt.Println(date)
	fmt.Println(strings.Repeat("=", 60))
	for _, period := range data {
		fmt.Println("["+period.MinutesString()+"]", period.Value)
	}
	stat := lib.Times{}
	stat = info.SumStat(stat)
	viewStatConsole(stat, info.Day)
}

// ---------------------------------------------------------------------------------------------------------------------
// Выводит массив статистики по категориям с заголовком
func viewStatConsole(stat lib.Times, title string) {
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println(title)
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
		fmt.Printf(warn()+" "+lib.FmtDuration(duration)+" ", category)
	}
	fmt.Printf("\n"+warn()+" ", "Total:")
	fmt.Println(lib.FmtDuration(total))
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
