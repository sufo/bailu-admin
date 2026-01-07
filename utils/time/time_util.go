/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 避免循环引用，增加了一层time文件夹
 */

package time

import (
	"fmt"
	"math"
	"github.com/sufo/bailu-admin/app/config"
	"net/http"
	"strconv"
	"time"
)

var (
	cst *time.Location
)

// CSTLayout China Standard Time Layout
const CSTLayout = "2006-01-02 15:04:05"

func init() {
	var err error
	if cst, err = time.LoadLocation(config.Conf.Server.TimeZone); err != nil {
		panic(err)
	}

	// 默认设置为中国时区
	time.Local = cst
}

// RFC3339ToCSTLayout convert rfc3339 value to china standard time layout
// 2020-11-08T08:18:46+08:00 => 2020-11-08 08:18:46
func RFC3339ToCSTLayout(value string) (string, error) {
	ts, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return "", err
	}

	return ts.In(cst).Format(CSTLayout), nil
}

// CSTLayoutString 格式化时间
// 返回 "2006-01-02 15:04:05" 格式的时间
func CSTLayoutString() string {
	ts := time.Now()
	return ts.In(cst).Format(CSTLayout)
}

func FormatCSTLayoutString(date time.Time) string {
	return date.In(cst).Format(CSTLayout)
}

// ParseCSTInLocation 格式化时间
func ParseCSTInLocation(date string) (time.Time, error) {
	return time.ParseInLocation(CSTLayout, date, cst)
}

// CSTLayoutStringToUnix 返回 unix 时间戳
// 2020-01-24 21:11:11 => 1579871471
func CSTLayoutStringToUnix(cstLayoutString string) (int64, error) {
	stamp, err := time.ParseInLocation(CSTLayout, cstLayoutString, cst)
	if err != nil {
		return 0, err
	}
	return stamp.Unix(), nil
}

// GMTLayoutString 格式化时间
// 返回 "Mon, 02 Jan 2006 15:04:05 GMT" 格式的时间
func GMTLayoutString() string {
	return time.Now().In(cst).Format(http.TimeFormat)
}

// ParseGMTInLocation 格式化时间
func ParseGMTInLocation(date string) (time.Time, error) {
	return time.ParseInLocation(http.TimeFormat, date, cst)
}

// SubInLocation 计算时间差
func SubInLocation(ts time.Time) float64 {
	return math.Abs(time.Now().In(cst).Sub(ts).Seconds())
}

func DaysBefore(day int) time.Time {
	current := time.Now()
	return current.AddDate(0, 0, day*(-1))
}

// is n days before
func IsNDaysAgo(day int, _date string) bool {
	daysAgo := DaysBefore(day)
	targetDate, err := time.Parse(CSTLayout, _date)
	return err == nil && targetDate.Before(daysAgo)
}

// 2024-5-2 19:00:00    0 0 19 2 5 ? 2024
func Time2CronExpression(t time.Time) string {
	////////////////////////// 秒 分	 时 日  月	年
	monthStr := t.Format("01")
	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s %s %s %d %d ? %d", t.Format("05"), t.Format("04"), t.Format("15"), t.Day(), month, t.Year())
}

// duration readable
func FormatDuration(d time.Duration) string {
	// 转换为更大的单位
	years := d / (365 * 24 * time.Hour)
	d %= 365 * 24 * time.Hour
	days := d / (24 * time.Hour)
	d %= 24 * time.Hour
	hours := d / time.Hour
	d %= time.Hour
	minutes := d / time.Minute

	// 构建格式化字符串
	parts := []string{}
	if years > 0 {
		parts = append(parts, fmt.Sprintf("%d年", years))
	}
	if days > 0 {
		parts = append(parts, fmt.Sprintf("%d天", days))
	}
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%d小时", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%d分钟", minutes))
	}

	// 如果duration小于1分钟，显示"小于1分钟"
	if len(parts) == 0 {
		return "小于1分钟"
	}

	return fmt.Sprint(parts)
}
