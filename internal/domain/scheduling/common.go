package scheduling

import (
	"regexp"
	"strings"
	"time"

	"github.com/chop1k/medods-test/internal/domain/models"
)

const dateFormat = "02-01-2006"

func trimTime(from time.Time) time.Time {
	return time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, from.Location())
}

func diffDays(from time.Time, to time.Time) int {
	diff := to.Sub(from)

	return int(diff.Hours() / 24)
}

func weekday(day string) time.Weekday {
	switch day {
	case "sunday":
		return time.Sunday
	case "monday":
		return time.Monday
	case "tuesday":
		return time.Tuesday
	case "wednesday":
		return time.Wednesday
	case "thursday":
		return time.Thursday
	case "friday":
		return time.Friday
	case "saturday":
		return time.Saturday
	default:
		return -1
	}
}

func weekNameToNumber(weeday string) int {
	var number int

	switch weeday {
	case "sunday":
		number = 1
	case "monday":
		number = 2
	case "tuesday":
		number = 3
	case "wednesday":
		number = 4
	case "thursday":
		number = 5
	case "frinday":
		number = 6
	case "saturday":
		number = 7
	}

	return number
}

func dateMatched(date, pattern string) bool {
	patternRegex := strings.ReplaceAll(regexp.QuoteMeta(pattern), `\*`, `\d`)

	matched, err := regexp.MatchString("^"+patternRegex+"$", date)
	if err != nil {
		return false
	}

	return matched
}

func scheduleWeeklyTask(templates models.Template, from time.Time, to time.Time) ([]models.TaskBody, error) {
	panic("not implemented")
}

func scheduleMonthlyTask(templates models.Template, from time.Time, to time.Time) ([]models.TaskBody, error) {
	panic("not implemented")
}

func scheduleOneshotTask(templates models.Template, from time.Time, to time.Time) ([]models.TaskBody, error) {
	panic("not implemented")
}
