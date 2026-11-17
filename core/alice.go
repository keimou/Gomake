package core

import (
	"gmkdb/global"
	"gmkdb/initalice"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// initialize project program
func init() {
	// 解决跨地区时区问题
	loc, _ := time.LoadLocation("Asia/Shanghai")
	time.Local = loc
	if false { // 开启生产环境 true || false 已替代 -> setMode()
		gin.SetMode(gin.ReleaseMode) // gin.DebugMode gin.ReleaseMode gin.TestMode
	}

	setMode() // command-line open production environment
	initialize()

	global.GIA_LOG.Info("Core System initialze completed.")
}

// start
func initialize() {
	Viper()             // initialize viper
	Zap()               // initialize zap operation log library
	initalice.Preload() // preload initialize modul
}

// set gin mode
func setMode() {
	// mode := osInput()
	var mode string
	if len(os.Args) > 1 {
		if os.Args[1] != "-c" {
			mode = os.Args[1]
		} else if len(os.Args) > 3 {
			mode = os.Args[3]
		}
	}

	switch mode {
	case "release", "0":
		// fmt.Printf("ReleaseMode: %v\n", gin.ReleaseMode)
		gin.SetMode(gin.ReleaseMode)
	case "test", "1":
		// fmt.Printf("TestMode: %v\n", gin.ReleaseMode)
		gin.SetMode(gin.TestMode)
	default: // default is DebugMode, no modification is required
		// fmt.Printf("DebugMode\n")
		// gin.SetMode(gin.DebugMode)
	}
}

// func osInput() (m string) {
// 	// fmt.Printf("os.Args[%d]: %v\n", len(os.Args), os.Args)
// 	if len(os.Args) > 1 {
// 		m = os.Args[1]
// 		if m == "-c" {
// 			if len(os.Args) < 4 {
// 				m = ""
// 			} else {
// 				m = os.Args[3]
// 			}
// 		}
// 	}
// 	return
// }
