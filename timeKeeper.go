package main

import (
	"fmt"
	"time"
)

var weekdayWords = [7]string{
	"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday",
}
var specialOrdinalWords = map[int]string{
	1: "1st", 2: "2nd", 3: "3rd", 21: "21st", 22: "22nd", 23: "23rd", 31: "31st",
}

func speakTimeEvents() {
	lastMinute := -1
	for {
		now := time.Now()
		thisMinute := now.Hour()*60 + now.Minute()
		if lastMinute != thisMinute {
			lastMinute = thisMinute
			if settings.Active && thisMinute%settings.IntervalMinutes == 0 {
				speak(englishTime(now))
			}
		}
		// pool somewhat often because host might've been asleep
		time.Sleep(time.Second)
	}
}

func englishTime(now time.Time) string {
	text := timeToEnglishText(now)
	return fmt.Sprintf("It is %s.", text)
}

func timeToEnglishText(now time.Time) string {
	hour, minute := now.Hour(), now.Minute()
	if minute != 0 {
		return fmt.Sprintf("%d %d", hour, minute)
	} else if hour == 12 {
		return "noon"
	} else if hour == 1 {
		return "1 hour"
	} else if hour != 0 {
		return fmt.Sprintf("%d hours", hour)
	}

	// it's midnight!
	weekday := weekdayWords[now.Weekday()]
	day := now.Day()
	dayWord := specialOrdinalWords[day]
	if dayWord != "" {
		return fmt.Sprintf("%s the %s", weekday, dayWord)
	} else {
		return fmt.Sprintf("%s the %dth", weekday, day)
	}
}
