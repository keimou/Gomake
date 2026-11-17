package global

import (
	"log"
	"time"
)

// cache system(csh) to clear
var TaskCsh interface{ Clear() } // 当前任务管理器缓存
var PsCsh interface{ Clear() }   // 发布订阅管理器缓存

// Cleanup system, done end, Shutdown
func Shutdown(start time.Time) {
	// ====== close code start ======
	st := time.Now()

	// // 追踪通道数据, 读完数据进行关闭.
	// // 避免关闭后任写入
	// go closeMainChan()

	// // close tasks
	// TaskCsh.Clear()
	// PsCsh.Clear()

	// // close the database link before the program ends
	// DBClose()
	// // 关闭 redis
	// for _, r := range rdbs {
	// 	r.Close()
	// }
	// Make sure to refresh the log when the program exits
	GIA_LOG.Sync()

	// ******* code end **********
	log.Println("resources are being cleaning...")
	elapsed := time.Since(start)
	sted := time.Since(st)
	log.Printf("cleaned completed, consumption %v\n  * system uptime: %v\n", sted, elapsed)
	// os.Exit(0) // Manual exit with exit code 0
}

// // 关闭主数据通道  // 放入任务清理, 先发生任务关闭, 再关闭通道然后再等待任务清理完成
// func closeMainChan() {
// 	if TaosDBCh != nil { // 确定已初始化
// 		close(TaosDBCh)
// 		// for {
// 		// 	n := len(TaosDBCh)
// 		// 	log.Printf("End Close MainChan len(%d)...\n", n)
// 		// 	time.Sleep(time.Duration(n/10) * time.Second)
// 		// 	if time.Sleep(1 * time.Second); n == 0 {
// 		// 		close(TaosDBCh)
// 		// 		break
// 		// 	}
// 		// }
// 		log.Println("End Close MainChan Completed.")
// 	}
// }
