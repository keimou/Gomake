package core

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"gmkdb/core/internal"
	"gmkdb/global"
)

func Viper(path ...string) *viper.Viper {
	var config string

	if len(path) == 0 {
		// 命令行 -c 参数输入
		flag.StringVar(&config, "c", "", "choose config file.")
		flag.Parse()
		// 未输入参数, 未传值, 系统默认指定文件
		if config == "" { // 判断命令行参数是否为空
			if configEnv := os.Getenv(internal.ConfigEnv); configEnv == "" { // 判断 internal.ConfigEnv 常量存储的环境变量是否为空
				switch gin.Mode() {
				case gin.DebugMode:
					config = internal.ConfigDefaultFile
					// fmt.Printf("您正在使用gin模式的%s环境名称,config的路径为%s\n", gin.EnvGinMode, internal.ConfigDefaultFile)
				case gin.ReleaseMode:
					config = internal.ConfigReleaseFile
					// fmt.Printf("您正在使用gin模式的%s环境名称,config的路径为%s\n", gin.EnvGinMode, internal.ConfigReleaseFile)
				case gin.TestMode:
					config = internal.ConfigTestFile
					// fmt.Printf("您正在使用gin模式的%s环境名称,config的路径为%s\n", gin.EnvGinMode, internal.ConfigTestFile)
				}
				log.Printf("正在使用gin模式的 %s 环境, config的路径为%s\n", gin.Mode(), config) // gin.EnvGinMode
			} else { // internal.ConfigEnv 常量存储的环境变量不为空 将值赋值于config
				config = configEnv
				fmt.Printf("正在使用 %s 环境变量, config的路径为%s\n", internal.ConfigEnv, config)
			}
		} else { // 命令行参数不为空 将值赋值于config
			fmt.Printf("正在使用命令行的-c参数传递的值,config的路径为%s\n", config)
		}
	} else { // 函数传递的可变参数的第一个值赋值于config
		config = path[0]
		fmt.Printf("正在使用func Viper()传递的值,config的路径为%s\n", config)
	}

	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		// err = fmt.Errorf("Fatal error config file: %s \n", err)
		panic(err)
	}
	v.WatchConfig()

	// 实时跟新系统配置
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err = v.Unmarshal(&global.GIA_CONFIG); err != nil {
			fmt.Println(err)
		}
		// initialize()
	})
	// 加载系统配置
	if err = v.Unmarshal(&global.GIA_CONFIG); err != nil {
		log.Println("GIA_CONFIG Loading Error:", err)
		// } else {
		// 	log.Println("config file is Ok!")
	}

	// fmt.Println("Viper is Ok!")
	// // root 适配性 根据root位置去找到对应迁移位置,保证root路径有效
	// global.GIA_CONFIG.AutoCode.Root, _ = filepath.Abs("..")
	// global.BlackCache = local_cache.NewCache(
	// 	local_cache.SetDefaultExpire(time.Second * time.Duration(global.GIA_CONFIG.JWT.ExpiresTime)),
	// )

	global.GIA_VPR = v
	return v
}
