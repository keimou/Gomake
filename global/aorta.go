package global

import (
	"gmkdb/config"

	"github.com/songzhibin97/gkit/cache/local_cache"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"golang.org/x/sync/singleflight"

	"gorm.io/gorm"

	"github.com/go-redis/redis/v8"
)

// global
// Golang Interface Architecture
var (
	GIA_VPR    *viper.Viper
	GIA_CONFIG config.Server

	GIA_DB     *gorm.DB
	GIA_DBList map[string]any // map[string]*gorm.DB
	GIA_REDIS  *redis.Client

	// Log
	GIA_LOG *zap.Logger
	Zap_gin *zap.Logger

	// // GVA_LOG    *oplogging.Logger
	// GVA_LOG                 *zap.Logger
	// GVA_Timer               timer.Timer = timer.NewTimerTask()

	// jwt
	// GIA_Concurrency_Control = &singleflight.Group{}
	JWT_Concurrency_Control = &singleflight.Group{} // jwt Token归并,避免并发
	BlackCache              local_cache.Cache       // jwt 黑名单缓存

	// lock sync.RWMutex
)

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
// 		panic("db no init")
// 	}
// 	return db
// }
