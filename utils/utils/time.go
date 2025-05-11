package utils

import (
	"math"
	"time"
)

const (
	ISO8601_UTC = "2006-01-02T15:04:05Z"
)

// 返回下周一中午12点的时间
func NextMondayMidnight() time.Time {
	today := time.Now()
	daysUntilMonday := (1 + 7 - int(today.Weekday())) % 7
	nextMonday := today.AddDate(0, 0, daysUntilMonday)
	nextMondayMidnight := time.Date(nextMonday.Year(), nextMonday.Month(), nextMonday.Day(), 0, 0, 0, 0, nextMonday.Location())
	return nextMondayMidnight
}

// 返回下个月1号0点的时间
func NextMonthFirstDay() time.Time {
	now := time.Now()
	nextMonth := now.AddDate(0, 1, 0)
	nextMonthFirstDay := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, nextMonth.Location())
	return nextMonthFirstDay
}

// 获取当前时间的0点
func GetMidnight() time.Time {
	now := time.Now()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return midnight
}

// 获取t的utc的0点
func GetMidnightWithUtc(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

// 比较时间;-1=time1早于/小于time2；1=time1晚于/大于time2；0=相同
func CompareTime(time1 time.Time, time2 time.Time) int {
	if time1.Before(time2) {
		return -1
	} else if time1.After(time2) {
		return 1
	} else {
		return 0
	}
}

// 计算两个时间直接差的天数,取绝对值
func DaysBetweenDates(date1, date2 time.Time) int64 {
	days := date2.Sub(date1).Hours() / 24
	return int64(math.Abs(days))
}

func GetDayBegin() time.Time {
	now := time.Now()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return midnight
}

// 获取Chart计算开始时间
func GetChartDayStart() time.Time {
	dayBegin := GetDayBegin()

	if dayBegin.Unix() < time.Now().Unix() {
		return dayBegin
	}

	return dayBegin.Add(time.Duration(-24) * time.Hour)
}

func GetChartWeekDayStart() time.Time {
	dayBegin := GetDayBegin()
	weekday := int(dayBegin.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	monday := dayBegin.AddDate(0, 0, 1-weekday)

	return monday
}

func GetCurrentTimeFormat() string {
	return time.Now().Format(time.DateTime)
}

// 获取本周周一的日期零点
func GetMondayFormatted(today time.Time) time.Time {
	weekday := today.Weekday()
	// 计算与周一的偏移量
	if weekday == 0 {
		weekday = 7
	}
	offset := int(-weekday + 1)

	// 计算周一的日期
	monday := today.AddDate(0, 0, offset)
	monday = time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, time.UTC)
	return monday
}

// GetTimeDifference 计算从当前时间到明天0点的时间差，以秒为单位
func GetTimeDifference() time.Duration {

	now := time.Now()
	// 明天的日期设置为0点
	tomorrowMidnight := GetDayBegin().AddDate(0, 0, 1)
	// 计算时间差
	diff := tomorrowMidnight.Sub(now)
	// 返回时间差的秒数
	return diff
}
