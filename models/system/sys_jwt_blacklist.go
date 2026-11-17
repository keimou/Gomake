package system

import "gmkdb/global"

// 存入数据库中的黑名单, 考虑放入redis做过期释放
type JwtBlacklist struct {
	// global.GIA_MODEL
	ID  uint   `gorm:"primarykey"`
	Jwt string `gorm:"type:text;comment:jwt"`
	global.TS_MODEL
}

// 到期时间字段: 用于清理数据库, 判断数据库内无效token
