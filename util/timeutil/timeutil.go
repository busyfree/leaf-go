package timeutil

import (
	"strconv"
	"strings"
	"time"
)

var timeLocation = time.Now().Location()
var Start time.Time
var End time.Time

func StartTick() {
	Start = time.Now()
}

func EndTick() float32 {
	delta := time.Now().Sub(Start)
	return float32(delta.Seconds())
}

func EndTickDt() time.Duration {
	return time.Now().Sub(Start)
}

func MsTimestampStr2Time(stampStr string) (time.Time, error) {
	stamp, err := strconv.ParseInt(stampStr, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return MsTimestamp2Time(stamp), nil
}

//毫秒转time对象
func MsTimestamp2Time(stamp int64) time.Time {
	return time.Unix(stamp/1000, 0)
}
func Timestamp2Time(stamp int64) time.Time {
	return time.Unix(stamp, 0)
}

//毫秒
func MsTimestampNow() int64 {
	return time.Now().UnixNano() / 1000000
}

// SecondTimestampNow 当前秒时间戳
func SecondTimestampNow() int64 {
	return time.Now().UnixNano() / 1000000000
}

func MsTimestamp2SecondStr() string {
	return time.Now().Format("20060102150405")
}

func MsTimestamp2DayStr() string {
	return time.Now().Format("20060102")
}

func MsTimestamp2DayHourStr() string {
	return time.Now().Format("2006010215")
}

func MsTimestamp2DayMiniuteStr() string {
	return time.Now().Format("200601021504")
}

func MsTimestamp2MonthStr() string {
	return time.Now().Format("200601")
}

func MsTimestamp2MilliStr() string {
	return strings.Replace(time.Now().Format("20060102150405.000"), ".", "", 1)
}

func MsTimestamp2MicroStr() string {
	return strings.Replace(time.Now().Format("20060102150405.000000"), ".", "", 1)
}
func MsTimestamp2NanoStr() string {
	return strings.Replace(time.Now().Format("20060102150405.000000000"), ".", "", 1)
}

func NanoTimestamp2Time(stamp int64) time.Time {
	return time.Unix(stamp/1000000000, stamp%1000000000)
}

func ParseTime(layout, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, timeLocation)
}

// this monday to next monday
func WeekRange(now time.Time) (weekStart, weekEnd time.Time) {
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	weekStart = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	weekEnd = weekStart.AddDate(0, 0, 7)
	return
}


func GetNextMonthRange(now time.Time, interval int) (firstOfMonth, firstOfNextMonth time.Time, lastDayWeek time.Weekday, lastDay int) {
	currentYear, currentMonth, _ := now.Date()
	firstOfMonth = time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, timeLocation)
	firstOfNextMonth = firstOfMonth.AddDate(0, interval, 0)
	lastDayTime := firstOfMonth.AddDate(0, interval, -1)
	lastDayWeek = lastDayTime.Weekday()
	lastDay = lastDayTime.Day()
	return
}

func GetBeforeMonthRange(now time.Time, interval int) (firstOfMonth, firstOfNextMonth time.Time, lastDayWeek time.Weekday, lastDay int) {
	currentYear, currentMonth, _ := now.Date()
	firstOfNextMonth = time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, timeLocation)
	firstOfMonth = firstOfNextMonth.AddDate(0, interval, 0)
	lastDayTime := firstOfNextMonth.AddDate(0, 0, -1)
	lastDayWeek = lastDayTime.Weekday()
	lastDay = lastDayTime.Day()
	return
}