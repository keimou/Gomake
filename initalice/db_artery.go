package initalice

import (
	"gmkdb/config"
	"gmkdb/global"
	"log"
	"os"

	"gorm.io/gorm"
)

// @Author: [keid]
// @Function: InitDB
// @Description: 初始化数据库连接
// @Date: 2025-02-28
func activateDataBase() {
	//  initialize primary database connect
	dbMain()
	// initialize other databases connect
	dbList()
	registerTables(global.GIA_DB) // 迁移系统数据库表
	if len(global.GIA_DBList) > 0 {
		for alias, db := range global.GIA_DBList {
			registerTables(db, alias) // 迁移从数据库表
		}
	}
	// initialize other server databases(X待整理)
	// bace.InitDB()
}

// Gorm 初始化数据库并产生数据库全局变量
// @Author: [keid]
// @Function: Gorm
// @Description: 初始化主数据库, 并产生数据库全局变量
// @Date: 2025-03-01
func dbMain() {
	// switch global.GIA_CONFIG.System.DbType {
	// case "mysql":
	// 	return GormMysql(config.Mysql{GeneralDB: global.GIA_CONFIG.DBMain})
	// default:
	// 	return nil
	// }

	// mainCfg := config.SpecializedDB{GeneralDB: global.GIA_CONFIG.DBMain, Type: global.GIA_CONFIG.System.DbType, AliasName: "MainDB"}
	if db, ok := openDB(*global.DBConfig()).(*gorm.DB); ok && db != nil {
		global.GIA_DB = db
	} else {
		// panic("Failed to connect Primary Database No *gorm.DB")
		// global.GIA_LOG.Panic("Failed to connect Primary Database")
		log.Println("No Primary Database")
		os.Exit(0)
	}
}

// @Author: [keid]
// @Function: DBList
// @Description: 初始化从数据库
// @Date: 2025-03-01
func dbList() {
	dbMap := make(map[string]any)
	for _, info := range global.GIA_CONFIG.DBList {
		if info.Disable {
			continue
		}
		// fmt.Printf("AliasName: %v - ", info.AliasName)
		if db := openDB(info); db != nil {
			dbMap[info.AliasName] = db
		}
	}

	// fmt.Printf("len(dbMap): %v\n", len(dbMap))
	global.GIA_DBList = dbMap
}
func openDB(info config.SpecializedDB) any {
	// fmt.Printf("Type: %v\n", info.Type)
	switch info.Type {
	case "sqlite":
		return OpenSqlite(info.GeneralDB)
	case "mysql":
		return OpenMysql(config.Mysql{GeneralDB: info.GeneralDB})
	// case "taos":
	// 	info.Port = "6041" // *玄学*: 这里如不用Restful连接,命令会阻塞无法执行,但程序中可用Sql链接不影响.  原因已找到,暂未处理,保留设置.
	// 	// 注册表结构的时候会将其替换回来
	// 	return global.TDengine(config.Taos{GeneralDB: info.GeneralDB})
	// 	// return nil
	default:
		log.Println("No Database Type for ", info.Type, info.Dbname, info.AliasName)
		os.Exit(0)
		return nil
	}
}
