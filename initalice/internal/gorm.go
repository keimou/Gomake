package internal

import (
	"log"
	"os"
	"time"

	"gmkdb/config"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBBASE interface {
	// GetLogMode() string
	GetLogMode() logger.LogLevel
	GetLogZap() bool
	GetLogLength() int
	GetBatchSize() int
	GetSlowThreshold() time.Duration
}

var Gorm = new(_gorm)

type _gorm struct{}

// @Author: [keid]
// @Function: Config
// @Description: Config gorm 自定义配置
// @Date: 2025-02-28
func (g *_gorm) Config(m interface{}) *gorm.Config {
	// config := &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}
	db_config := &gorm.Config{}
	// CreateBatchSize: 500,

	var modeCfg DBBASE
	switch v := m.(type) {
	case config.Mysql:
		modeCfg = &v
	// case config.Taos:
	// 	modeCfg = &v
	default:
		return db_config // 默认 silent
	}

	logLevel := modeCfg.GetLogMode()

	_default := logger.New(NewWriter(log.New(os.Stdout, "\r\n", log.LstdFlags), logLevel, modeCfg.GetLogZap(), modeCfg.GetLogLength()), logger.Config{
		SlowThreshold: modeCfg.GetSlowThreshold(),
		LogLevel:      logLevel, // logger.Warn
		Colorful:      true,
		// IgnoreRecordNotFoundError: false,
		// ParameterizedQueries:      false,
	})
	// db_config.CreateBatchSize = modeCfg.GetBatchSize() // if is 0 set off for default
	num := modeCfg.GetBatchSize()
	if num > 0 {
		db_config.CreateBatchSize = num
	}
	db_config.Logger = _default.LogMode(logLevel)

	return db_config
}

// ----------

// // Config gorm 自定义配置
// // Author [SliverHorn](https://github.com/SliverHorn)
// func (g *_gorm) Config() *gorm.Config {
// 	config := &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}
// 	_default := logger.New(NewWriter(log.New(os.Stdout, "\r\n", log.LstdFlags)), logger.Config{
// 		SlowThreshold: 200 * time.Millisecond,
// 		LogLevel:      logger.Warn,
// 		Colorful:      true,
// 	})
// 	var logMode DBBASE
// 	switch global.GIA_CONFIG.System.DbType {
// 	case "mysql":
// 		logMode = &global.GIA_CONFIG.Mysql
// 		break
// 	// case "pgsql":
// 	// 	logMode = &global.GIA_CONFIG.Pgsql
// 	// 	break
// 	default:
// 		logMode = &global.GIA_CONFIG.Mysql
// 	}

// 	switch logMode.GetLogMode() {
// 	case "silent", "Silent":
// 		config.Logger = _default.LogMode(logger.Silent)
// 	case "error", "Error":
// 		config.Logger = _default.LogMode(logger.Error)
// 	case "warn", "Warn":
// 		config.Logger = _default.LogMode(logger.Warn)
// 	case "info", "Info":
// 		config.Logger = _default.LogMode(logger.Info)
// 	default:
// 		config.Logger = _default.LogMode(logger.Info)
// 	}
// 	return config
// }
