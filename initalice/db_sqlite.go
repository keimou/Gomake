package initalice

import (
	"gmkdb/config"
	"gmkdb/global"
	"gmkdb/initalice/internal"

	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	// _ "modernc.org/sqlite" // 可不用开启GOGCC
)

func OpenSqlite(s config.GeneralDB) *gorm.DB {
	if s.Dbname == "" {
		return nil
	}

	dbConfig := sqlite.Config{
		DSN: s.Dbname + ".sqlite",
	}

	if db, err := gorm.Open(sqlite.New(dbConfig), internal.Gorm.Config(s)); err != nil {
		global.GIA_LOG.Panic("Sqlite connect failed for"+s.Dbname, zap.Error(err))
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(s.MaxIdleConns)
		sqlDB.SetMaxOpenConns(s.MaxOpenConns)

		if err := sqlDB.Ping(); err != nil {
			global.GIA_LOG.Panic("Sqlite Ping failed for "+s.Dbname, zap.Error(err))
		} else {
			global.GIA_LOG.Debug("Sqlite connection is successful", zap.String(s.Dbname, "PONG"))
		}
		return db
	}
}
