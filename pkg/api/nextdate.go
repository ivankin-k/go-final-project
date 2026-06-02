package api

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

// Compare dates
func after(a, b time.Time) bool {
	if a.Truncate(24 * time.Hour).After(b.Truncate(24 * time.Hour)) {
		return true
	}
	return false
}

func nextDateYear(now, dstart time.Time) time.Time {
	for {
		dstart = dstart.AddDate(1, 0, 0)
		if after(dstart, now) {
			break
		}
	}
	return dstart
}

// Function to convert:
// [1,5,-1,-2] to [1,5,30,31] or [1,5,29,30]
// [10, 31]	   to [10,31]     or [10]
func getRealDatesOfMonth(year, month int, days []int) []int {
	var (
		dates     []int
		numOfDays int
	)

	numOfDays = time.Date(year, time.Month(month+1), 0, 0, 0, 0, 0, time.UTC).Day()

	for _, day := range days {
		if day > 0 && day <= numOfDays {
			dates = append(dates, day)
		}
		if day < 0 {
			dates = append(dates, numOfDays+1+day)
		}
	}

	slices.Sort(dates)
	return dates
}

func nextDateMonth(now, dstart time.Time, repeatParams []string) (time.Time, error) {

	// Find next date starting from a bigger date, as no need to move by steps (as in 'd' or 'y')
	if after(dstart, now) {
		now = dstart
	}

	// Check if any params
	if len(repeatParams) == 0 {
		return time.Time{}, fmt.Errorf("Missing parameters for 'm' repeat rule")
	}

	// Parse/validate taskDaysRule list
	var (
		err          error
		day          int
		taskDaysRule []int
	)
	for _, val := range strings.Split(repeatParams[0], ",") {
		if day, err = strconv.Atoi(val); err != nil {
			return time.Time{}, fmt.Errorf("Can't convert 'm' repeat rule day param to int: %s\n%v\n", val, err)
		}
		if day < -2 || day > 31 || day == 0 {
			return time.Time{}, fmt.Errorf("Invalid 'm' repeat rule day parameter: %d (valid: 1-31, -1, -2)", day)
		}
		taskDaysRule = append(taskDaysRule, day)
	}

	// Parse/validate taskMonths list (if any)
	var (
		month      int
		taskMonths []int
	)
	if len(repeatParams) > 1 {
		for _, val := range strings.Split(repeatParams[1], ",") {
			if month, err = strconv.Atoi(val); err != nil {
				return time.Time{}, fmt.Errorf("Can't convert 'm' repeat rule month param to int: %s\n%v\n", val, err)
			}
			if month < 1 || month > 12 {
				return time.Time{}, fmt.Errorf("Invalid 'm' repeat rule month parameter: %d (valid: 1-12)", month)
			}
			taskMonths = append(taskMonths, month)
		}
		slices.Sort(taskMonths)
	} else {
		taskMonths = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	}

	var (
		nowDay,
		nowMonth,
		nowYear int

		taskDays []int
	)
	nowYear, _, nowDay = now.Date()
	nowMonth = int(now.Month())

	// Get real dates available within current month
	taskDays = getRealDatesOfMonth(nowYear, nowMonth, taskDaysRule)

	// If nothing is available within current year, set the first available date next year
	if (taskMonths[len(taskMonths)-1] <= nowMonth) && (taskDays[len(taskDays)-1] <= nowDay) {
		return time.Date(nowYear+1, time.Month(taskMonths[0]), taskDays[0], 0, 0, 0, 0, time.UTC), nil
	}

	// If dates are available in current month
	if slices.Contains(taskMonths, nowMonth) {
		for _, day := range taskDays {
			if day > nowDay {
				return time.Date(nowYear, now.Month(), day, 0, 0, 0, 0, time.UTC), nil
			}
		}
	}

	// Return next available month's first date
	for _, month := range taskMonths {
		if month > nowMonth {
			// Get dates within next available month
			taskDays = getRealDatesOfMonth(nowYear, month, taskDaysRule)
			return time.Date(nowYear, time.Month(month), taskDays[0], 0, 0, 0, 0, time.UTC), nil
		}
	}

	return time.Time{}, nil
}

func nextDateWeek(now, dstart time.Time, repeatParams []string) (time.Time, error) {

	// Find next date starting from a bigger date, as no need to move by steps (as in 'd' or 'y')
	if after(dstart, now) {
		now = dstart
	}

	// Check if any params
	if len(repeatParams) == 0 {
		return time.Time{}, fmt.Errorf("Missing parameters for 'w' repeat rule")
	}

	var (
		err      error
		weekDays []int
		weekDay  int
	)

	// Parse/validate weekdays in rule
	for _, v := range strings.Split(repeatParams[0], ",") {
		if weekDay, err = strconv.Atoi(v); err != nil {
			return time.Time{}, fmt.Errorf("Can't convert 'w' repeat rule weekday param to int: %s\n%v\n", v, err)
		}
		if weekDay < 1 || weekDay > 7 {
			return time.Time{}, fmt.Errorf("Invalid 'w' repeat rule weekday parameter: %d (valid: 1-7)", weekDay)
		}
		weekDays = append(weekDays, weekDay)
	}
	slices.Sort(weekDays)

	var taskDate time.Time
	taskDate = now

	// Find next date
	for {
		taskDate = taskDate.AddDate(0, 0, 1)
		// convert Sun-Sat to Mon-Sun week and check
		if slices.Contains(weekDays, ((int(taskDate.Weekday())+6)%7)+1) {
			break
		}
	}

	return taskDate, nil
}

func nextDateDay(now, dstart time.Time, repeatParams []string) (time.Time, error) {

	// Validate days parameter
	if len(repeatParams) == 0 {
		return time.Time{}, fmt.Errorf("Missing parameter for 'd' repeat rule")
	}
	var (
		err     error
		addDays int
	)
	if addDays, err = strconv.Atoi(repeatParams[0]); err != nil {
		return time.Time{}, fmt.Errorf("Can't convert 'd' repeat rule param to int: %s\n%v\n", repeatParams[0], err)
	}
	if addDays < 0 || addDays > 400 {
		return time.Time{}, fmt.Errorf("Invalid 'd' repeat rule parameter: %d (valid: 1-400)", addDays)
	}

	// Find next date
	for {
		dstart = dstart.AddDate(0, 0, addDays)
		if after(dstart, now) {
			break
		}
	}

	return dstart, nil
}

func nextDate(now time.Time, dstart string, repeat string) (string, error) {

	// Validate parameters
	if len(dstart) == 0 {
		return "", fmt.Errorf("'date' parameter is missing")
	}
	var (
		err   error
		dtime time.Time
	)
	if dtime, err = time.Parse(dateFormat, dstart); err != nil {
		return "", fmt.Errorf("Can't parse %q to date (YYYYMMDD):\n%v\n", dstart, err)
	}
	if len(repeat) < 1 {
		return "", fmt.Errorf("Repeat rule is undefined")
	}

	// Split repeat rule and route to related funcs
	repeatParams := strings.Split(repeat, " ")
	switch repeatParams[0] {
	case "y":
		dtime = nextDateYear(now, dtime)
	case "m":
		if dtime, err = nextDateMonth(now, dtime, repeatParams[1:]); err != nil {
			return "", err
		}
	case "w":
		if dtime, err = nextDateWeek(now, dtime, repeatParams[1:]); err != nil {
			return "", err
		}
	case "d":
		if dtime, err = nextDateDay(now, dtime, repeatParams[1:]); err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("Invalid task repeat rule: %q", repeatParams[0])
	}

	return dtime.Format(dateFormat), nil

}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	var (
		now    time.Time
		result string
		err    error
	)

	// Check if 'now' is received
	if now, err = time.Parse(dateFormat, r.FormValue("now")); err != nil {
		now = time.Now()
	}

	// Calculate next date
	if result, err = nextDate(now, r.FormValue("date"), r.FormValue("repeat")); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Send response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}
