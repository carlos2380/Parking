package utils

import (
	"log"
	"strconv"
	"time"
)

func CentsIntToEurStr(cents int) string {

	euros := float64(cents) / 100
	eurosStr := strconv.FormatFloat(euros, 'f', 2, 64)
	return eurosStr
}

func DateToStr(date time.Time) string {
	return date.Format(time.RFC822)
}

func StrDateToDate(dateStr string) (time.Time, error) {
	date, err := time.Parse(time.RFC822, dateStr)
	if err != nil {
		log.Fatalf("Error on parse entry date %s : %v", dateStr, err)
		return time.Time{}, err
	}
	return date, nil
}

func GetMinutes(startTime *time.Time, endTime *time.Time) int {
	duration := endTime.Sub(*startTime)
	return int(duration.Minutes())
}
