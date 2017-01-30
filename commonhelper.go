package main

import (
	"time"
	"strconv"
	"strings"
)

func getWeekNumber() string{
	_, week := time.Now().ISOWeek()
	weekStr := strconv.Itoa(week)
	return weekStr;
}

func GetWeekDays() [7]string{
	year, week := time.Now().ISOWeek()
	var weekArray [7]string
	firstDay := FirstDayOfISOWeek(year, week, time.FixedZone("UTC", 3)) // Minsk time zone


	for i := 0; i < 7; i++ {
		nextDay:= firstDay.AddDate(0, 0, i)
		weekArray[i] = nextDay.Format("2006-01-2")
	}
/*	formatedFirstDay := firstDay.Format("2006-01-2")
	nextday:=firstDay.AddDate(0, 0, 1)
	fmt.Print(nextday.Format("2006-01-2"))
	weekArray := [2]string {formatedFirstDay, nextday}*/
	return weekArray
}

func GetWeekDaysForWeekNumber(week int ) [7]string{
	year, _ := time.Now().ISOWeek()
	var weekArray [7]string
	firstDay := FirstDayOfISOWeek(year, week, time.FixedZone("UTC", 3)) // Minsk time zone


	for i := 0; i < 7; i++ {
		nextDay:= firstDay.AddDate(0, 0, i)
		weekArray[i] = nextDay.Format("2006-01-2")
	}
	/*	formatedFirstDay := firstDay.Format("2006-01-2")
		nextday:=firstDay.AddDate(0, 0, 1)
		fmt.Print(nextday.Format("2006-01-2"))
		weekArray := [2]string {formatedFirstDay, nextday}*/
	return weekArray
}

func ExtendStringSlice(slice []string, element string) []string {
	n := len(slice)
	if n == cap(slice) {
		// Slice is full; must grow.
		// We double its size and add 1, so if the size is zero we still grow.
		newSlice := make([]string, len(slice), 2*len(slice)+1)
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0 : n+1]
	slice[n] = element
	return slice
}

func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}


func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func FirstDayOfISOWeek(year int, week int, timezone *time.Location) time.Time {
	date := time.Date(year, 0, 0, 0, 0, 0, 0, timezone)
	isoYear, isoWeek := date.ISOWeek()

	// iterate back to Monday
	for date.Weekday() != time.Monday {
		date = date.AddDate(0, 0, -1)
		isoYear, isoWeek = date.ISOWeek()
	}

	// iterate forward to the first day of the first week
	for isoYear < year {
		date = date.AddDate(0, 0, 7)
		isoYear, isoWeek = date.ISOWeek()
	}

	// iterate forward to the first day of the given week
	for isoWeek < week {
		date = date.AddDate(0, 0, 7)
		isoYear, isoWeek = date.ISOWeek()
	}

	return date
}


