package main

import (
	"gmkdb/core"
)

// @title Swagger Example API
// @version 0.0.1
// @description This is a API Server
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name sbl-token
// @BasePath /
func main() {
	// global.GIA_VPR = core.Viper() //初始化viper
	// // // 初始化zap日志库
	// global.GIA_LOG = core.Zap() // 初始化zap日志库

	// global.GIA_DB = initalice.Gorm() // 初始化数据库连接
	// // // 初始化计时器

	// // 初始化数据列表
	// initalice.DBList()

	// if global.GIA_DB != nil {
	// 	initalice.RegisterTables(global.GIA_DB) // 初始化表
	// 	db, _ := global.GIA_DB.DB()
	// 	defer db.Close()
	// }

	core.RunServer()
}
