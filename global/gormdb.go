package global

import (
	"database/sql"
	"fmt"
	"gmkdb/config"
	"sync"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var lock sync.RWMutex

// // GetGlobalDBByDBName 通过名称获取db list中的db
// func GetGlobalDBByDBName(dbname string) *gorm.DB {
// 	lock.RLock()
// 	defer lock.RUnlock()
// 	return GIA_DBList[dbname]
// }

// // MustGetGlobalDBByDBName 通过名称获取db 如果不存在则panic
// func MustGetGlobalDBByDBName(dbname string) *gorm.DB {
// 	lock.RLock()
// 	defer lock.RUnlock()
// 	db, ok := GIA_DBList[dbname]
// 	if !ok || db == nil {
// 		// panic("db no init")
//		os.Exit(0)
// 	}
// 	return db
// }

func DBListGorm(aliasName string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	if db, ok := GIA_DBList[aliasName].(*gorm.DB); ok {
		return db
	}
	return nil
}

func DBListSql(aliasName string) *sql.DB {
	lock.RLock()
	defer lock.RUnlock()
	if db, ok := GIA_DBList[aliasName].(*sql.DB); ok {
		return db
	}
	return nil
}

// @Author: [keid]
// @Function: DBClose
// @Description: 关闭数据库集连接,在程序结束之前关闭所有数据库连接
// @Date: 2025-02-28
func DBClose() {
	// close the database link before the program ends
	if len(GIA_DBList) > 0 {
		for alias, d := range GIA_DBList {
			fmt.Println(" - 关闭数据库连接:", alias)
			dbClose(d)
		}
	}
	dbClose(GIA_DB)
	GIA_LOG.Debug("End Close Database Completed.")
}

// 关闭数据库连接(检查连接是否为空且未关闭)
func dbClose(db any) {
	switch d := db.(type) {
	case *gorm.DB:
		if d != nil {
			sqlDB, _ := d.DB()
			if sqlDB.Ping() == nil {
				sqlDB.Close()
			}
		}
	case *sql.DB:
		if d != nil {
			if d.Ping() == nil {
				d.Close()
			}
		}
	case *redis.Client:
		fmt.Printf("\n\nEnd dbClose: redis - %v\n", d)
	default:
		fmt.Printf("\n\nEnd dbClose: %v\n", d)
	}
}

// @Author: [keid]
// @Function: DBConfig
// @Description: 获取数据库配置信息
// @Date: 2025-02-28
//
//	若不传参数, 则返回主数据库配置
//	主数据库由system中db-aliias指定, 若未指定或指定无效结束程序 (废弃, 已使用固定系统主数据库设置)
func DBConfig(alias ...string) (cfg *config.SpecializedDB) {
	// dbName := GIA_CONFIG.System.DbAlias
	if len(alias) > 0 {
		aliasName := alias[0]
		for _, info := range GIA_CONFIG.DBList {
			if info.Disable {
				continue
			} else if info.AliasName == aliasName {
				return &info
			}
		}
	} else {
		return &config.SpecializedDB{
			Disable:   false,
			Type:      GIA_CONFIG.System.DbType,
			AliasName: "System MainDB", // System Primary DB
			GeneralDB: GIA_CONFIG.DBMain,
		}
	}
	return
}
