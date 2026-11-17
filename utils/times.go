package utils

import (
	"fmt"
	"regexp"
	"time"
)

// 将 Interval 转换为对应的 time (1Moon=30day)
func TimeInterval(interval string) (tsGap time.Duration) {
	number, letter := StringSubmatch(interval)
	switch letter {
	case "M":
		// tsAfter = tsTime.Truncate(time.Minute).AddDate(0, number, 0).UnixMilli()
		tsGap = time.Duration(number) * 30 * 24 * time.Hour
	case "w":
		tsGap = time.Duration(number) * 7 * 24 * time.Hour
	case "d":
		tsGap = time.Duration(number) * 24 * time.Hour
	case "h":
		tsGap = time.Duration(number) * time.Hour
	case "m":
		tsGap = time.Duration(number) * time.Minute
	case "s":
		tsGap = time.Duration(number) * time.Second
	default:
	}
	return
}

// 计算时间 A,B 的时间差
func TimeDiff(startTime, endTime int64) time.Duration {
	// 定义两个时间戳
	timeA := time.UnixMilli(startTime)
	timeB := time.UnixMilli(endTime)

	// 计算时间差
	duration := timeB.Sub(timeA) // B-A

	return duration
}

func TimeDiffStr(startTime, endTime int64) string {
	// 计算时间差
	duration := TimeDiff(startTime, endTime)

	var timeDiff string
	if duration < 0 {
		timeDiff += "-"
	}
	if int(duration.Hours()) != 0 {
		timeDiff += fmt.Sprintf("%dh", int(duration.Abs().Hours()))
	}
	if int(duration.Minutes())%60 != 0 {
		timeDiff += fmt.Sprintf("%02dm", int(duration.Abs().Minutes())%60)
	}
	if int(duration.Seconds())%60 != 0 {
		timeDiff += fmt.Sprintf("%02ds", int(duration.Abs().Seconds())%60)
	}
	if (timeDiff == "" || timeDiff == "-") && duration != 0 {
		return duration.String()
	}

	return timeDiff
}

// int64 毫秒时间戳 按分钟时间格式输出 Millisecond
func TransTimeMin(timeMs int64) string {
	seconds := timeMs / 1000
	nanoseconds := (timeMs % 1000) * int64(time.Millisecond)
	// 将时间戳转换为 time.Time 对象
	ts := time.Unix(seconds, nanoseconds)
	return ts.Format("2006-01-02 15:04")
}

// int64 毫秒时间戳 按秒时间格式输出 Millisecond
func TransTimeSec(timeMs int64, fs ...int) string {
	seconds := timeMs / 1000
	nanoseconds := (timeMs % 1000) * int64(time.Millisecond)
	// 将时间戳转换为 time.Time 对象
	ts := time.Unix(seconds, nanoseconds)

	tsStr := ts.Format("2006-01-02 15:04:05")
	if len(fs) > 0 {
		switch fs[0] {
		case 0:
			tsStr = ts.Format("15:04:05")
		case 1:
			tsStr = ts.Format("15:04")
		case 2:
			tsStr = ts.Format("01-02 15:04")
		case 3:
			tsStr = ts.Format("01-02 15:04:05")
		}
	}

	return tsStr
}

// int64 毫秒时间戳 按毫秒时间格式输出 Millisecond
func TransTime(timeMs int64) string {
	seconds := timeMs / 1000
	nanoseconds := (timeMs % 1000) * int64(time.Millisecond)
	// 将时间戳转换为 time.Time 对象
	ts := time.Unix(seconds, nanoseconds)
	return ts.Format("2006-01-02 15:04:05.000")
}

// 时间间隔, 向后推一个间隔时间(未来)
func IntervalAfter(ts int64, interval string, multiplier ...int) (tsAfter int64) {
	tsTime := time.UnixMilli(ts)
	// time.UnixMilli(ts).Add(time.Minute).UnixMilli()
	// fmt.Printf("transTime: %v\n", transTime)
	number, letter := StringSubmatch(interval)
	if len(multiplier) > 0 && multiplier[0] > 1 {
		number *= multiplier[0]
	}
	switch letter {
	case "M":
		tsAfter = tsTime.AddDate(0, number, 0).UnixMilli()
	case "w":
		tsGap := time.Duration(number) * 7 * 24 * time.Hour
		tsAfter = tsTime.Add(tsGap).UnixMilli()
	case "d":
		tsGap := time.Duration(number) * 24 * time.Hour
		tsAfter = tsTime.Add(tsGap).UnixMilli()
	case "h":
		tsGap := time.Duration(number) * time.Hour
		tsAfter = tsTime.Add(tsGap).UnixMilli()
	case "m": // .Truncate(time.Minute)
		tsGap := time.Duration(number) * time.Minute
		tsAfter = tsTime.Add(tsGap).UnixMilli()
	case "s": // .Truncate(tsGap)
		tsGap := time.Duration(number) * time.Second
		tsAfter = tsTime.Add(tsGap).UnixMilli()
	default:
		tsAfter = ts
	}

	return tsAfter
}

// 时间间隔, 向前推一个间隔时间(历史)
func IntervalBefore(ts int64, interval string, multiplier ...int) (tsBefore int64) {
	tsTime := time.UnixMilli(ts)
	number, letter := StringSubmatch(interval)
	if len(multiplier) > 0 && multiplier[0] > 1 {
		number *= multiplier[0]
	}
	switch letter {
	case "M":
		tsBefore = tsTime.AddDate(0, -number, 0).UnixMilli()
	case "w":
		tsBefore = tsTime.AddDate(0, 0, -number*7).UnixMilli()
	case "d":
		tsBefore = tsTime.AddDate(0, 0, -number).UnixMilli()
	case "h":
		tsGap := time.Duration(number) * time.Hour
		tsBefore = tsTime.Add(-tsGap).UnixMilli()
	case "m":
		tsGap := time.Duration(number) * time.Minute
		tsBefore = tsTime.Add(-tsGap).UnixMilli() // .Truncate(time.Minute)
	case "s":
		tsGap := time.Duration(number) * time.Second
		tsBefore = tsTime.Add(-tsGap).UnixMilli() // .Truncate(time.Second)
	default:
		tsBefore = ts
	}
	return tsBefore
}

//	interval切片: 将start_time 的 mid_time 转换为 long_time; 如传入的是当前时间，则需向前取1m .Add(-time.Minute)???

// 取'nowTime'对应的'start_time': 通过传入时间校对当前'interval'取整获取'start_time'
func IntervalOf(nowTime int64, interval string) (startTime int64) {

	// 将时间戳转换为时间对象
	tsTime := time.UnixMilli(nowTime)

	number, letter := StringSubmatch(interval)
	switch letter {
	case "s":
		startTime = tsTime.Truncate(time.Duration(number) * time.Second).UnixMilli()
	case "m":
		startTime = tsTime.Truncate(time.Duration(number) * time.Minute).UnixMilli()
	case "h":
		startTime = tsTime.Truncate(time.Duration(number) * time.Hour).UnixMilli()
	case "d":
		startTime = tsTime.Truncate(time.Duration(number) * 24 * time.Hour).UnixMilli()
	case "w":
		startTime = tsTime.Truncate(time.Duration(number) * 7 * 24 * time.Hour).UnixMilli()
	case "M":
		// startTime = time.Date(tsTime.Year(), tsTime.Month(), 1, 0, 0, 0, 0, tsTime.Location()).UnixMilli()
		startTime = time.Date(tsTime.Year(), tsTime.Month(), 1, 0, 0, 0, 0, time.UTC).UnixMilli()
	}
	// var gapTs int64
	// gapTs = time.Minute.Milliseconds() * int64(number)
	// longTime = (nowTime / gapTs) * gapTs

	return
}

// 取'Time'对应的取整时间向后跳一间隔
func IntervalGap(nowTime int64, tsGap time.Duration) time.Time {
	tsTime := time.UnixMilli(nowTime)
	return tsTime.Truncate(tsGap).Add(tsGap)
}

// startTime 用方法计算
// endTime 为下一个startTime-1
// 查询时
// 用 map[string]struct{} 存放所有时间线
// for map 1m直接存 3m-1M 计算当前数据的 startTime ,比较v

// @Author: [keid]
// @Function: ToTime
// @Description: 文本格式转化时间格式
// @Date: 2026-10-28
func ToTime(timeStr string) (time.Time, error) {
	// Go 的参考时间格式，必须严格按照这个写法
	// 2006(年) 01(月) 02(日) 15(时) 04(分) 05(秒)
	layout := generateLayout(timeStr)

	// 关键：使用 ParseInLocation 并指定时区
	// 如果你的字符串代表的是本地时间（如北京时间），用 time.Local
	// 如果是 UTC 时间，用 time.UTC
	// t, err := time.ParseInLocation(layout, timeStr, time.Local)
	// if err != nil {
	// 	log.Println("时间解析失败:", err)
	// }

	// return time.ParseInLocation(layout, timeStr, time.FixedZone("CST", 8*3600))
	return time.ParseInLocation(layout, timeStr, time.Local)
}

// 提取标准时
func generateLayout(timeStr string) string {
	// 定义 Go time 包的标准 layout 数字占位符（按年月日时分秒顺序）
	layoutDigits := []string{"2006", "01", "02", "15", "04", "05"}
	index := 0

	re := regexp.MustCompile(`\d+`)
	result := re.ReplaceAllStringFunc(timeStr, func(match string) string {
		if index >= len(layoutDigits) {
			return match // 如果数字段多于预期，保留原样
		}
		replacement := layoutDigits[index]
		index++
		return replacement
	})

	return result
}
