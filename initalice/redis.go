package initalice

import (
	"context"
	"gmkdb/global"
	"log"

	"github.com/go-redis/redis/v8"
)

func Redis() {
	redisCfg := global.GIA_CONFIG.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password, // no password set
		DB:       redisCfg.DB,       // use default DB
	})
	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		// global.GIA_LOG.Error("redis connect ping failed, err:", zap.Error(err))
		log.Println("redis connect ping failed, err:", err.Error())
	} else {
		// global.GIA_LOG.Info("redis connect ping response:", zap.String("pong", pong))
		log.Println("redis connect ping response:", pong)
		global.GIA_REDIS = client
	}
}
