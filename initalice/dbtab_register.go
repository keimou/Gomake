package initalice

import (
	"gmkdb/global"
	"gmkdb/models/ninex"
	"gmkdb/models/system"

	// "gmkdb/daos/daosql"
	// "gmkdb/models/bnbmod"
	"gmkdb/utils"
	"os"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 注册数据库表

// @Author: [keid]
// @Function: RegisterTables
// @Description: 注册数据库表
func registerTables(dbValue any, alias ...string) {
	var aliasName string
	if len(alias) > 0 {
		aliasName = alias[0]
	}
	switch db := dbValue.(type) {
	case *gorm.DB:
		if utils.IsNull(aliasName) {
			registerPRI(db)
		} else {
			registerGorm(db, aliasName)
		}
		// case *sql.DB:
		// 	registerSql(db, aliasName)
	}
	// if len(dbName) > 0 {
	// 	registerGorm(db, dbName[0])
	// } else {
	// 	registerPRI(db)
	// }
}

// @Author: [keid]
// @Function: registerPRI
// @Description: 1)注册主数据库表专用
// @Date: 2025-02-27
func registerPRI(db *gorm.DB) {
	err := db.AutoMigrate(
		ninex.NexExpert{},
		ninex.NexRecord{},
		// 系统模块表
		// system.SysUser{},
		// system.SysAuthority{},
		// system.SysBaseMenu{},
		// system.SysAuthorityMenu{},
		// system.JwtBlacklist{},
		// system.SysOperationRecord{},
	// 扩展模块
	// example.ExaFileUploadAndDownload{},
	)
	if err != nil {
		global.GIA_LOG.Error("System MainDB("+global.GIA_CONFIG.System.DbType+") Register Table Failed.", zap.Error(err))
		os.Exit(0)
	}
	global.GIA_LOG.Debug("System MainDB(" + global.GIA_CONFIG.System.DbType + ") Register Table Success.")
}

// @Author: [keid]
// @Function: registerGorm
// @Description: 2)注册从数据库表 Replica Database
// @Date: 2025-02-28
func registerGorm(db *gorm.DB, aliasName string) {
	var err error
	switch aliasName {
	case "prime":
		err = db.AutoMigrate(
			// 系统模块表
			// system.SysUser{},
			// system.SysAuthority{},
			// system.SysBaseMenu{},
			// system.SysAuthorityMenu{},
			system.JwtBlacklist{},
			system.SysOperationRecord{},
		)
	case "bacdb":
		err = db.AutoMigrate(
		// // 基础模块
		)
	default:
		global.GIA_LOG.Error("["+aliasName+"] no set databse.", zap.String("gorm", aliasName))
		return
	}

	if err != nil {
		global.GIA_LOG.Error("["+aliasName+"] register table failed.", zap.Error(err))
		os.Exit(0)
	}
	global.GIA_LOG.Debug("[" + aliasName + "] register table success.")
}

// // @Author: [keid]
// // @Function: registerGorm
// // @Description: 2)注册从数据库表 Replica Database (taos -> sqldb)
// // @Date: 2025-02-28
// func registerSql(db *sql.DB, aliasName string) {
// 	var err error

// 	dbCfg := global.DBConfig(aliasName)
// 	dao := daosql.NewDao(db, dbCfg.Dbname)
// 	// // 结束后close掉, 记得踢出DBList
// 	// delete(global.GIA_DBList, aliasName) // 这里踢出会影响循环
// 	// dao.Close()

// 	if aliasName == "bactd" {
// 		err = dao.AutoMigrate(
// 		// bnbmod.KlineEvent{},
// 		)
// 	} else if aliasName == "qptd" {
// 		err = dao.AutoMigrate(
// 		// bnbmod.Rsi{},
// 		// bnbmod.EMA{Period: 50},
// 		// bnbmod.EMA{Period: 12}, // 26 12-144-169-576-676
// 		// bnbmod.EMA{Period: 144},
// 		// bnbmod.EMA{Period: 169},
// 		// bnbmod.EMA{Period: 576},
// 		// bnbmod.EMA{Period: 676},
// 		// bnbmod.MACD{},
// 		// bnbmod.Boll{},
// 		// bnbmod.KDJ{},
// 		// bnbmod.OBV{},
// 		// bnbmod.AO{},
// 		// bnbmod.ADX{},
// 		)
// 	} else {
// 		global.GIA_LOG.Error("["+aliasName+"] no set databse.", zap.String("sql", aliasName))
// 		return
// 	}

// 	if err != nil {
// 		global.GIA_LOG.Error("["+aliasName+"] register table failed.", zap.Error(err))
// 		os.Exit(0)
// 	}

// 	global.GIA_LOG.Debug("[" + aliasName + "] register table success.")
// 	if db.Ping() == nil {
// 		// fmt.Printf("taosCf: %v '%v' '%v'\n", dbCfg.Taos(), db, db == nil)
// 		if dbCfg.Taos().Port != "6041" { // 替换回原始连接
// 			db.Close()
// 			global.GIA_DBList[dbCfg.AliasName] = global.TDengine(dbCfg.Taos())
// 		}
// 		// fmt.Println("dbCfg: ", dbCfg, global.GIA_DBList[dbCfg.AliasName].(*sql.DB), db, db == nil)
// 	}
// }

// 系统如需要初始化数据,可进行操作
// 场景设置: 可识别系统中配置文件数据库参数是否存在,若存在直接读取,若不存在则初始化数据库
