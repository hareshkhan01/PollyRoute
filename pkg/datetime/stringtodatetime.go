package datetime

import "time"

func StringToDateTime(dateTimeString string) (time.Time, error) {

	// Parse the date-time string into a time.Time object
	parsedTime, err := time.Parse(time.RFC3339, dateTimeString)
	if err != nil {
		return time.Time{}, err
	}
	return parsedTime, nil
}
