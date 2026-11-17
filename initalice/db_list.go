package initalice

// import (
// 	"fmt"
// 	"gmkdb/config"
// 	"gmkdb/global"

// 	"gorm.io/gorm"
// )

// // const sys = "system"

// // func DBList() {
// // 	fmt.Println("Join DBlist!")
// // 	dbMap := make(map[string]*gorm.DB)
// // 	for _, info := range global.GIA_CONFIG.DBList {
// // 		if info.Disable {
// // 			continue
// // 		}
// // 		switch info.Type {
// // 		case "mysql":
// // 			dbMap[info.AliasName] = GormMysqlByConfig(config.Mysql{GeneralDB: info.GeneralDB})
// // 		// case "pgsql":
// // 		// 	dbMap[info.AliasName] = GormPgSqlByConfig(config.Pgsql{GeneralDB: info.GeneralDB})
// // 		default:
// // 			continue
// // 		}
// // 	}
// // 	// 做特殊判断,是否有迁移
// // 	// 适配低版本迁移多数据库版本
// // 	if sysDB, ok := dbMap[sys]; ok {
// // 		global.GIA_DB = sysDB
// // 	}
// // 	global.GIA_DBList = dbMap
// // }
