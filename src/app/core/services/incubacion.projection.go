package services

import (
	"fmt"
	"strconv"
	"time"
)

type Incubation struct {
	Fecha       time.Time
	WeekDay     time.Weekday
	Nacimientos int64
}

func GetIncubations(year string) []Incubation {
	defer bench("GetIncubations")
	_year, err := strconv.ParseInt(year, 10, 64)
	handleErr(err)
	_, rows := AbuelosProjectionTable(year, "nac", true)
	var mon_acc int64 = 0
	var fri_acc int64 = 0
	date, err := time.Parse("2006-03-01", fmt.Sprintf("%d-01-01", _year))
	handleErr(err)
	results := []Incubation{}
	for m := 1; m < len(rows); m++ {
		for d := 1; d < len(rows[0]); d++ {
			if int64(date.Year()) != _year {
				continue
			}
			if date.Weekday() > time.Monday && date.Weekday() <= time.Friday {
				fri_acc += rows[date.Month()][date.Day()]
				if date.Weekday() == time.Friday {
					el := Incubation{
						Fecha:       date,
						WeekDay:     date.Weekday(),
						Nacimientos: fri_acc,
					}
					results = append(results, el)
					fri_acc = 0
				}
			} else {
				mon_acc += rows[date.Month()][date.Day()]
				if date.Weekday() == time.Monday {
					el := Incubation{
						Fecha:       date,
						WeekDay:     date.Weekday(),
						Nacimientos: mon_acc,
					}
					results = append(results, el)
					mon_acc = 0
				}
			}
			date = date.AddDate(0, 0, 1)
		}
	}
	return results
}

func LisIncubations(year string) {
	incubations := GetIncubations(year)
	days := map[time.Weekday]string{
		1: "Lunes",
		5: "Viernes",
	}
	for _, inc := range incubations {
		fmt.Printf("%v - %s: %d\n", inc.Fecha.String(), days[inc.WeekDay], inc.Nacimientos)
	}
}
