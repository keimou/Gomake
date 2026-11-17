package daords

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
)

type RdsDao struct {
	rdb *redis.Client
}

// # INSERT DELETE SELECT UPDATE
func NewDao(rdb *redis.Client) *RdsDao {
	return &RdsDao{
		rdb: rdb,
	}
}

var once sync.Once

func (r *RdsDao) Abc() {
	once.Do(func() {
		pong, err := r.rdb.Ping(context.Background()).Result()
		if err != nil {
			fmt.Printf(" - redis无法链接 err: %v\n", err)
		} else {
			fmt.Printf(" * redis功能为开发 ping:%v (%v)\n", pong, r.rdb)
		}
	})
}
