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

func emitTimeEvents(timeEvents chan time.Time) {
	lastMinute := -1
	for {
		now := time.Now()
		thisMinute := now.Hour()*60 + now.Minute()
		if lastMinute != thisMinute {
			lastMinute = thisMinute
			if settings.Active && thisMinute%settings.IntervalMinutes == 0 {
				timeEvents <- now
			}
		}
		// pool somewhat often in order to still be useful when host is suspended
		time.Sleep(time.Second)
	}
}

func emitEnglishSpeechEvents(speakEvents chan string) {
	timeEvents := make(chan time.Time)
	go emitTimeEvents(timeEvents)
	for now := range timeEvents {
		speakEvents <- timeToEnglishText(now)
	}
}

func timeToEnglishText(now time.Time) string {
	hour, minute := now.Hour(), now.Minute()
	minuteText := "hours"
	if minute == 0 {
		if hour == 0 {
			weekday := weekdayWords[now.Weekday()]
			day := specialOrdinalWords[now.Day()]
			if day == "" {
				day = fmt.Sprintf("%dth", now.Day())
			}
			return fmt.Sprintf("It is %s the %s.", weekday, day)
		} else if hour == 12 {
			return "It is noon."
		}
		if hour == 1 {
			minuteText = "hour"
		}
	} else {
		minuteText = fmt.Sprint(now.Minute())
	}
	return fmt.Sprintf("It is %d %s.", hour, minuteText)
}
