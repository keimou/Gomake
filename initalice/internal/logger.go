package internal

import (
	"fmt"
	"gmkdb/global"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"gorm.io/gorm/logger"
)

type writer struct {
	logger.Writer
	logLevel  logger.LogLevel // 添加日志级别字段
	logZap    bool            // 是否启用 zap
	logLength int
}

// NewWriter writer 构造函数
// @Author: [keid]
// @Function: NewWriter
// @Description: 重构 logger.Writer 函数, 自定义输出
// @Date: 2025-03-01
func NewWriter(w logger.Writer, logLevel logger.LogLevel, logZap bool, logLength int) *writer {
	return &writer{
		Writer:    w,
		logLevel:  logLevel,
		logZap:    logZap,
		logLength: logLength,
	}
}

// @Author: [keid]
// @Function: Printf
// @Description: 格式化打印gorm数据库日志
// @Date: 2025-03-01
func (w *writer) Printf(message string, data ...any) {
	if len(data) > 0 {
		// fmt.Printf("len(top): %v\n", TotalLength(data))
		// 按 logLength 切割数据信息
		data = processData(w.logLength, data)
		// fmt.Printf("len(end): %v\n", TotalLength(data))
	}
	logMessage := fmt.Sprintf("[GORM] "+message+"\n", data...)
	if w.logZap {
		switch w.logLevel {
		case logger.Error:
			global.GIA_LOG.Error(logMessage)
		case logger.Warn:
			global.GIA_LOG.Warn(logMessage)
		case logger.Info:
			global.GIA_LOG.Info(logMessage)
		default:
			global.GIA_LOG.Debug(logMessage)
		}
	} else {
		// w.Writer.Printf(message, data...)
		wMessage := fmt.Sprintf(message, data...)
		// timestamp := time.Now().Format("2006-01-02 15:04:05")
		wMessage = fmt.Sprintf("-\t %s  %s", getLevel(w.logLevel), wMessage)
		w.Writer.Printf("%s", wMessage)
	}
}

// 过长数据截取对应长度
func processData(maxLength int, data []any) []any {
	if maxLength == 0 { // 判断 data 是否有值 || len(data) == 0
		return data
	}
	var result []any
	currentLength := 0
	for _, item := range data {
		var str string
		switch v := item.(type) {
		case string:
			str = v
		case fmt.Stringer:
			str = v.String()
		case float64:
			str = strconv.FormatFloat(v, 'f', 3, 64)
		case int64:
			str = strconv.FormatInt(v, 10)
		default:
			str = fmt.Sprintf("%v", v)
		}
		remaining := maxLength - currentLength // 计算剩余需要的长度
		if len(str) > remaining {              // 剩余长度不足, 截取字符串
			str = str[:remaining] + "..."
		}
		// 添加到结果中
		result = append(result, convType(str))
		currentLength += len(str)
		if currentLength >= maxLength {
			return result
		}
	}
	return data
}

// 原始输出格式级别获取
func getLevel(logLevel logger.LogLevel) string {
	switch logLevel {
	case logger.Error:
		return color.New(color.BgRed, color.FgHiWhite).SprintFunc()("  Error ")
	case logger.Warn:
		return color.New(color.BgHiYellow).SprintFunc()("  Warn  ")
	case logger.Info:
		return color.New(color.BgHiGreen).SprintFunc()("  Info  ")
	default:
		return color.New(color.FgHiBlue).SprintFunc()("  Debug ")
	}
}

func convType(s string) any {
	// Remove spaces before and after strings
	s = strings.TrimSpace(s)
	// Try to resolve to int64
	if n, err := strconv.ParseInt(s, 10, 64); err == nil {
		return n
	}
	// Try to resolve to float64
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	// If it's neither, return string
	return s
}

// 若float没有小数位, 则补0
// if !strings.Contains(str, ".") { // strconv.FormatFloat(v, 'f', -1, 64)
// str += ".00"
// }

// // 计算总数据长度
// func TotalLength(data []interface{}) int {
// 	total := 0
// 	for _, item := range data {
// 		var str string
// 		switch v := item.(type) {
// 		case string:
// 			str = v
// 		case fmt.Stringer:
// 			str = v.String()
// 		default:
// 			str = fmt.Sprintf("%v", v)
// 		}
// 		total += len(str)
// 	}
// 	return total
// }

// -----------

// type writer struct {
// 	logger.Writer
// }

// // NewWriter writer 构造函数
// // Author [SliverHorn](https://github.com/SliverHorn)
// func NewWriter(w logger.Writer) *writer {
// 	return &writer{Writer: w}
// }

// // Printf 格式化打印日志
// // Author [SliverHorn](https://github.com/SliverHorn)
// func (w *writer) Printf(message string, data ...interface{}) {
// 	var logZap bool
// 	switch global.GIA_CONFIG.System.DbType {
// 	case "mysql":
// 		logZap = global.GIA_CONFIG.Mysql.LogZap
// 		// case "pgsql":
// 		// 	logZap = global.GIA_CONFIG.Pgsql.LogZap
// 	}
// 	if logZap {
// 		// 	global.GIA_LOG.Info(fmt.Sprintf(message+"\n", data...))
// 		// } else {
// 		w.Writer.Printf(message, data...)
// 	}
// }
