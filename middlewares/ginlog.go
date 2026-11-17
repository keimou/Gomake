package middlewares

import (
	"fmt"
	"gmkdb/global"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GinLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var strColor *color.Color
		start := time.Now()

		// 处理请求
		c.Next()

		encodeLevel := global.GIA_CONFIG.Zap.EncodeLevel
		showColor := true
		switch encodeLevel {
		case "CapitalLevelEncoder", "LowercaseLevelEncoder":
			showColor = false
		}

		// 记录日志
		// 使用颜色库: White Black Red Green Yellow Blue Cyan Magenta
		// 请求状态码
		status := c.Writer.Status()
		staStr := fmt.Sprintf(" %3v ", status)
		if showColor {
			switch {
			case status >= 500 || status < 100:
				strColor = color.New(color.BgRed, color.FgHiWhite)
			case status >= 400:
				strColor = color.New(color.BgHiYellow)
			case status >= 300 || status < 200:
				strColor = color.New(color.BgHiWhite)
			default:
				strColor = color.New(color.BgHiGreen) // status >= 200
			}
			staStr = strColor.SprintFunc()(staStr)
		}
		// 请求方法
		method := c.Request.Method
		methodStr := fmt.Sprintf(" %-8v", method)
		if showColor {
			switch method {
			case "GET":
				strColor = color.New(color.BgBlue, color.FgHiWhite)
			case "POST":
				strColor = color.New(color.BgHiCyan)
			case "PUT":
				strColor = color.New(color.BgHiYellow)
			case "DELETE":
				strColor = color.New(color.BgRed, color.FgHiWhite)
			case "PATCH":
				strColor = color.New(color.BgHiGreen)
			case "OPTIONS":
				strColor = color.New(color.BgHiWhite)
			case "HEAD":
				strColor = color.New(color.BgMagenta, color.FgHiWhite)
			default: // CONNECT TRACE
				strColor = color.New()
			}
			methodStr = strColor.SprintFunc()(methodStr)
		}
		// 请求路径参数
		path := c.Request.URL.Path
		if c.Request.URL.RawQuery != "" {
			path += "?" + c.Request.URL.RawQuery
		}

		logger.Info(fmt.Sprintf("|%v|%14v |%16v |%v `%v`", staStr, time.Since(start), c.ClientIP(), methodStr, path))
		// if status == 404 {
		// 	logger.Debug(fmt.Sprintf("|%v|%14v |%16v |%v `%v`", staStr, time.Since(start), c.ClientIP(), methodStr, path))
		// 	// } else if status >= 300 {
		// 	// 	logger.Warn(fmt.Sprintf("|%v|%14v |%16v |%v `%v`", staStr, time.Since(start), c.ClientIP(), methodStr, path))
		// } else {
		// 	logger.Info(fmt.Sprintf("|%v|%14v |%16v |%v `%v`", staStr, time.Since(start), c.ClientIP(), methodStr, path))
		// }
	}
}

// [GIN] 2025/03/03 - 17:02:29 | 200 |     382.107µs |     192.168.0.1 | GET      "/health"
// [crypto]2025/03/03 - 17:04:06.227       INFO    /home/yun/Project/CryptoDB/initalice/internal/ginlog.go:26      /health {"status": 200, "method": "GET", "path": "/health", "query": "", "ip": "192.168.0.1", "user-agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36 Edg/133.0.0.0", "latency": 0.000060001}
// [GIN] 2025/03/03 - 17:04:06 |status(5)|latency(14) |ip(16) | method(9-1) "path'?'query"

func GinRecovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 检查是否是网络错误
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// 记录堆栈信息
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					c.Error(err.(error)) // 将错误传递给 Gin
					c.Abort()
					return
				}

				// 记录 panic 信息
				logger.Error("[Recovery from panic]",
					zap.Time("time", time.Now()),
					zap.Any("error", err),
					zap.String("request", string(httpRequest)),
					zap.String("stack", string(debug.Stack())),
				)

				// 返回 500 错误
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
