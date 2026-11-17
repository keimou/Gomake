package initalice

func Preload() {
	activateDataBase() // Initialize database connect
	// attachApply()      // Link Apply module for Exit
	activateChan(1000) // Activate data channel of taos
}

// Activate
// Attach

// // 运行结束后清理系统资源, 资源定向
// func attachApply() {
// 	global.TaskCsh = task.MGT
// 	global.PsCsh = pubsub.MGT
// }

func activateChan(n int) {
	// global.TaosDBCh = make(chan any, n)
}
