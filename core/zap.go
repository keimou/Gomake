package core

import (
	"fmt"
	"gmkdb/core/internal"
	"gmkdb/global"
	"gmkdb/utils"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Zap 获取 zap.Logger
// Author [SliverHorn](https://github.com/SliverHorn)
func Zap() (logger *zap.Logger) {
	if ok, _ := utils.PathExists(global.GIA_CONFIG.Zap.Director); !ok { // 判断是否有Director文件夹
		fmt.Printf("create %v directory\n", global.GIA_CONFIG.Zap.Director)
		_ = os.Mkdir(global.GIA_CONFIG.Zap.Director, os.ModePerm)
	}

	cores := internal.Zap.GetZapCores()
	logger = zap.New(zapcore.NewTee(cores...))
	// zap.NewProduction()
	// zap.NewNop()

	if global.GIA_CONFIG.Zap.ReplaceGin {
		global.Zap_gin = logger
	}

	if global.GIA_CONFIG.Zap.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}

	global.GIA_LOG = logger
	zap.ReplaceGlobals(global.GIA_LOG)
	return
}
