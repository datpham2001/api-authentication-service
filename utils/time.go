package utils

import "time"

func GetCurrentTimeZoneVN() time.Time {
	t := time.Now()

	location, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		panic(err)
	}

	return t.In(location)
}
