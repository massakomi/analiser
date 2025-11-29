package console

import (
	"analiser/pkg/lib"
	"fmt"
	"math"
	"time"
)

func printTest() {
	data := lib.WeekStatSorted()
	maxDuration := time.Duration(0)
	for _, week := range data {
		if week.Duration > maxDuration {
			maxDuration = week.Duration
		}
	}
	for _, week := range data {
		fmt.Println(math.Round(100 * (week.Duration.Seconds() / maxDuration.Seconds())))
	}
}
