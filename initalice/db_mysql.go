package initalice

import (
	"gmkdb/config"
	"gmkdb/global"
	"gmkdb/initalice/internal"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// @Author: [keid]
// @Function: GormMysqlByConfig
// @Description: 通过传入配置信息, 初始化Mysql数据库
// @Date: 2025-02-26
func OpenMysql(m config.Mysql) *gorm.DB {
	if m.Dbname == "" {
		return nil
	}
	mysqlConfig := mysql.Config{
		DSN:                       m.Dsn(), // DSN data source name
		DefaultStringSize:         191,     // string 类型字段的默认长度
		SkipInitializeWithVersion: false,   // 根据版本自动配置
	}

	if db, err := gorm.Open(mysql.New(mysqlConfig), internal.Gorm.Config(m)); err != nil {
		global.GIA_LOG.Panic("Mysql connect failed for "+m.Dbname, zap.Error(err))
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		// 验证连接
		if err := sqlDB.Ping(); err != nil {
			global.GIA_LOG.Panic("Mysql Ping failed for "+m.Dbname, zap.Error(err))
		} else {
			global.GIA_LOG.Debug(" Mysql connection is successful.", zap.String(m.Dbname, "PONG"))
		}
		return db
	}
}

// func isMySqlExist(db *sql.DB, dbName string) (ok bool, err error) {
// 	rows, err := db.Query("SHOW DATABASE LIKE ?", dbName)
// 	if err != nil {
// 		return
// 	}
// 	defer rows.Close()

// 	if rows.Next() {
// 		return true, nil
// 	}
// 	return false, rows.Err()
// }

// -----------

// // GormMysql 初始化Mysql数据库
// // Author [piexlmax](https://github.com/piexlmax)
// // Author [SliverHorn](https://github.com/SliverHorn)
// func GormMysql() *gorm.DB {
// 	fmt.Println("Join Mysql")

// 	m := global.GIA_CONFIG.Mysql
// 	if m.Dbname == "" {
// 		return nil
// 	}
// 	mysqlConfig := mysql.Config{
// 		DSN:                       m.Dsn(), // DSN data source name
// 		DefaultStringSize:         191,     // string 类型字段的默认长度
// 		SkipInitializeWithVersion: false,   // 根据版本自动配置
// 	}
// 	if db, err := gorm.Open(mysql.New(mysqlConfig), internal.Gorm.Config()); err != nil {
// 		return nil
// 	} else {
// 		sqlDB, _ := db.DB()
// 		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
// 		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
// 		return db
// 	}
// }

// // GormMysqlByConfig 初始化Mysql数据库用过传入配置
// func GormMysqlByConfig(m config.Mysql) *gorm.DB {
// 	if m.Dbname == "" {
// 		return nil
// 	}
// 	mysqlConfig := mysql.Config{
// 		DSN:                       m.Dsn(), // DSN data source name
// 		DefaultStringSize:         191,     // string 类型字段的默认长度
// 		SkipInitializeWithVersion: false,   // 根据版本自动配置
// 	}
// 	if db, err := gorm.Open(mysql.New(mysqlConfig), internal.Gorm.Config()); err != nil {
// 		panic(err)
// 	} else {
// 		sqlDB, _ := db.DB()
// 		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
// 		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
// 		return db
// 	}
// }
