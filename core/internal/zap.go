package internal

import (
	"fmt"
	"gmkdb/global"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Zap = new(_zap)

type _zap struct {
	prefix string
}

// GetEncoder 获取 zapcore.Encoder 2.编码器
// Author [SliverHorn](https://github.com/SliverHorn)
func (z *_zap) GetEncoder() zapcore.Encoder {
	if global.GIA_CONFIG.Zap.Format == "json" {
		return zapcore.NewJSONEncoder(z.GetEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(z.GetEncoderConfig())
}

// GetEncoderConfig 获取zapcore.EncoderConfig 编码器设置
// Author [SliverHorn](https://github.com/SliverHorn)
func (z *_zap) GetEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  global.GIA_CONFIG.Zap.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    global.GIA_CONFIG.Zap.ZapEncodeLevel(),
		EncodeTime:     z.CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder, // SecondsDurationEncoder(便于解析计算) StringDurationEncoder(易读)
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 完整路径:FullCallerEncoder 相对路径:ShortCallerEncoder
	}
}

// GetZapCores 根据配置文件的Level获取 []zapcore.Core 入口
// Author [SliverHorn](https://github.com/SliverHorn)
func (z *_zap) GetZapCores() []zapcore.Core {
	cores := make([]zapcore.Core, 0, 7)
	for level := global.GIA_CONFIG.Zap.TransportLevel(); level <= zapcore.FatalLevel; level++ {
		cores = append(cores, z.GetEncoderCore(level, z.GetLevelPriority(level)))
	}
	// 添加核心 (gin) info error
	return cores
}

// GetEncoderCore 获取Encoder的 zapcore.Core 1.核心
// Author [SliverHorn](https://github.com/SliverHorn)
func (z *_zap) GetEncoderCore(l zapcore.Level, levelFunc zap.LevelEnablerFunc, prefixs ...string) zapcore.Core {
	if len(prefixs) > 0 {
		z.prefix = prefixs[0]
	} else {
		z.prefix = global.GIA_CONFIG.Zap.Prefix
	}
	// [GIN]: no "level", INFO ERROR

	writer, err := FileRotatelogs.GetWriteSyncer(l.String()) // 使用file-rotatelogs进行日志分割
	if err != nil {
		fmt.Printf("Get Write Syncer Failed err:%v", err.Error())
		return nil
	}

	// 编码器, 文件流, 日志级别
	return zapcore.NewCore(z.GetEncoder(), writer, levelFunc)
}

// CustomTimeEncoder 自定义日志输出时间格式
// Author [SliverHorn](https://github.com/SliverHorn)
func (z *_zap) CustomTimeEncoder(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	// encoder.AppendString(t.Format(global.GEA_CONFIG.Zap.Prefix + " 2006/01/02 - 15:04:05.000"))
	encoder.AppendString(t.Format(z.prefix + " 2006/01/02 - 15:04:05.000"))
}

// GetLevelPriority 根据 zapcore.Level 获取 zap.LevelEnablerFunc
// Author [SliverHorn](https://github.com/SliverHorn)
func (z *_zap) GetLevelPriority(level zapcore.Level) zap.LevelEnablerFunc {
	switch level {
	case zapcore.DebugLevel:
		return func(level zapcore.Level) bool { // 调试级别 -1
			return level == zap.DebugLevel
		}
	case zapcore.InfoLevel:
		return func(level zapcore.Level) bool { // 日志级别 0
			return level == zap.InfoLevel
		}
	case zapcore.WarnLevel:
		return func(level zapcore.Level) bool { // 警告级别 1
			return level == zap.WarnLevel
		}
	case zapcore.ErrorLevel:
		return func(level zapcore.Level) bool { // 错误级别 2
			return level == zap.ErrorLevel
		}
	case zapcore.DPanicLevel:
		return func(level zapcore.Level) bool { // dpanic级别 3
			return level == zap.DPanicLevel
		}
	case zapcore.PanicLevel:
		return func(level zapcore.Level) bool { // panic级别 4
			return level == zap.PanicLevel
		}
	case zapcore.FatalLevel:
		return func(level zapcore.Level) bool { // 终止级别 5
			return level == zap.FatalLevel
		}
	default:
		return func(level zapcore.Level) bool { // 调试级别
			return level == zap.DebugLevel
		}
	}
}
