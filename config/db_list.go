package config

import (
	"time"

	"gorm.io/gorm/logger"
)

type DsnProvider interface {
	Dsn() string
}

// Embeded 结构体可以压平到上一层，从而保持 config 文件的结构和原来一样
// 见 playground: https://go.dev/play/p/KIcuhqEoxmY

// GeneralDB 也被 Pgsql 和 Mysql 原样使用
type GeneralDB struct {
	Path         string `mapstructure:"path" json:"path" yaml:"path"`                               // 服务器地址:端口
	Port         string `mapstructure:"port" json:"port" yaml:"port"`                               //:端口
	Config       string `mapstructure:"config" json:"config" yaml:"config"`                         // 高级配置
	Dbname       string `mapstructure:"db-name" json:"db-name" yaml:"db-name"`                      // 数据库名
	Username     string `mapstructure:"username" json:"username" yaml:"username"`                   // 数据库用户名
	Password     string `mapstructure:"password" json:"password" yaml:"password"`                   // 数据库密码
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"max-idle-conns" yaml:"max-idle-conns"` // 空闲中的最大连接数
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"max-open-conns" yaml:"max-open-conns"` // 打开到数据库的最大连接数
	LogMode      string `mapstructure:"log-mode" json:"log-mode" yaml:"log-mode"`                   // 是否开启Gorm全局日志
	LogLength    int    `mapstructure:"log-length" json:"log-length" yaml:"log-length"`             // 单条日志数据最大长度
	LogZap       bool   `mapstructure:"log-zap" json:"log-zap" yaml:"log-zap"`                      // 是否通过zap写入日志文件

	SlowThreshold int `mapstructure:"slow-threshold" json:"slow-threshold" yaml:"slow-threshold"` // 数据库慢速阈值-多少毫秒执行每批次
	BatchSize     int `mapstructure:"batch-size" json:"batch-size" yaml:"batch-size"`             // 数据库分批存储最大数据量
}

type SpecializedDB struct {
	Disable   bool   `mapstructure:"disable" json:"disable" yaml:"disable"`
	Type      string `mapstructure:"type" json:"type" yaml:"type"`
	AliasName string `mapstructure:"alias-name" json:"alias-name" yaml:"alias-name"`
	GeneralDB `yaml:",inline" mapstructure:",squash"`
}

//	func (m *GeneralDB) GetLogMode() string {
//		return m.LogMode
//	}

func (m *GeneralDB) GetBatchSize() int {
	return m.BatchSize
}

func (m *GeneralDB) GetLogMode() logger.LogLevel {
	switch m.LogMode {
	case "silent", "Silent", "SILENT":
		return logger.Silent
	case "error", "Error", "ERROR":
		return logger.Error
	case "warn", "Warn", "WARN":
		return logger.Warn
	case "info", "Info", "INFO":
		return logger.Info
	default:
		return logger.Info
	}
}

func (m *GeneralDB) GetLogZap() bool {
	return m.LogZap
}

func (m *GeneralDB) GetLogLength() int {
	return m.LogLength
}

func (m *GeneralDB) GetSlowThreshold() time.Duration {
	return time.Duration(m.SlowThreshold) * time.Millisecond
}
