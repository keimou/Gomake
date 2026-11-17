package initalice

import (
	"gmkdb/global"
	"gmkdb/middlewares"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/webdav"

	ginSwagger "github.com/swaggo/gin-swagger"
	// "github.com/swaggo/gin-swagger/swaggerFiles"

	_ "gmkdb/docs"
)

func Routers() *gin.Engine {
	Router := initRouter()

	// Router.StaticFS(global.GIA_CONFIG.Local.Path, http.Dir(global.GIA_CONFIG.Local.StorePath)) // 为用户头像和文件提供静态地址
	// Router.StaticFile("/favicon.ico", "/favicon.ico")
	// // Router.Use(middleware.LoadTls())  // 如果需要使用https 请打开此中间件
	// // // 然后前往 core/server.go 将启动模式 更变为 Router.RunTLS("端口","你的cre/pem文件","你的key文件")

	// // 跨域，如需跨域可以打开下面的注释
	// Router.Use(middlewares.Cors())        // 直接放行全部跨域请求
	Router.Use(middlewares.CorsByRules()) // 按照配置的规则放行跨域请求
	global.GIA_LOG.Debug("use middleware cors;")

	// Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler)) // gin-swagger@1.4->1.6: swaggerFiles.Handler = &webdav.Handler{FileSystem: webdav.NewMemFS(), LockSystem: webdav.NewMemLS()}
	Router.GET("/swagger/*any", ginSwagger.WrapHandler(&webdav.Handler{})) // FileSystem: webdav.NewMemFS(), LockSystem: webdav.NewMemLS()
	global.GIA_LOG.Debug("register swagger handler;")
	// 方便统一 添加路由组前缀 多服务器上线使用

	// baceRouter := routers.RouterGroupApp.Bace
	// applyRouter := routers.RouterGroupApp.Apply
	// qtzRouter := routers.RouterGroupApp.Qtz
	// systemRouter := routers.RouterGroupApp.System
	// exampleRouter := routers.RouterGroupApp.Example

	PublicGroup := Router.Group("")
	{
		// 健康监测
		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(200, "ok")
		})
		if false {
			initdetectRouter(PublicGroup)
		}
		// systemRouter.InitBaseRouter(PublicGroup) // 注册基础功能路由 不做鉴权
	}

	PrivateGroup := Router.Group("binance") // 私有组
	PrivateGroup.Use()                      // PrivateGroup.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	{
		// v1Api := PrivateGroup.Group("api/v1")
		// 注册 binance 接口模块
		// baceRouter.InitWsKlineRouter(v1Api)
		// baceRouter.InitKlineRouter(v1Api)
		// baceRouter.InitSymbolInfoRouter(v1Api) // Symbols info module

		// qtzRouter.InitQPRouter(v1Api)
		// qtzRouter.InitEMARouter(v1Api)
		// qtzRouter.InitRsiRouter(v1Api)
		// qtzRouter.InitMacdRouter(v1Api)
		// qtzRouter.InitBollRouter(v1Api)
		// qtzRouter.InitKDJRouter(v1Api)
		// qtzRouter.InitOBVRouter(v1Api)
		// qtzRouter.InitAORouter(v1Api)

		// 系统应用接口
		// applyRouter.InitTaskRouter(v1Api) // 注册系统任务功能
		// 注册系统异步通讯功能
		// systemRouter.InitApiRouter(PrivateGroup)                // 注册功能api路由
		// systemRouter.InitJwtRouter(PrivateGroup)                // jwt相关路由
		// systemRouter.InitUserRouter(PrivateGroup)               // 注册用户路由
		// systemRouter.InitMenuRouter(PrivateGroup)               // 注册menu路由
		// systemRouter.InitAuthorityRouter(PrivateGroup)          // 注册角色路由
		// systemRouter.InitCasbinRouter(PrivateGroup)             // 权限相关路由
		// systemRouter.InitSysOperationRecordRouter(PrivateGroup) // 操作记录
		// systemRouter.InitSystemRouter(PrivateGroup)             // system相关路由
	}

	global.GIA_LOG.Debug("Router register success.")
	return Router
}

func initRouter() (engine *gin.Engine) {
	if global.GIA_CONFIG.Zap.ReplaceGin {
		engine = gin.New()
		// // 使用自定义日志中间件
		engine.Use(middlewares.GinLogger(global.Zap_gin))
		engine.Use(middlewares.GinRecovery(global.Zap_gin))
	} else {
		engine = gin.Default()
	}
	initSvcApply()
	return engine
}

// 初始化功能应用
func initSvcApply() {

	// // 提前设置好消费者
	// w := services.ServiceGroupApp.BaceSvcGroup.WsKlineService
	// // 开启 Consumer 数据通道消费者 (默认3)
	// if err := w.Consumer(2, 1); err != nil {
	// 	global.GIA_LOG.Warn("Router Open default consumer Failed.", zap.Error(err))
	// }

	// // 链接 wskline 数据流
	// if err := w.WsKlineStream(); err != nil {
	// 	global.GIA_LOG.Warn("Router Link kline stream Failed.", zap.Error(err))
	// }

	// 开启一个定时器,
	automated()
}

func automated() {

}

// Automated Programs
// automation

//  -----------

// func Routers() *gin.Engine {
// 	Router := gin.Default()
// 	systemRouter := routers.RouterGroupApp.System
// 	// exampleRouter := routers.RouterGroupApp.Example

// 	// 如果想要不使用nginx代理前端网页，可以修改 web/.env.production 下的
// 	// VUE_APP_BASE_API = /
// 	// VUE_APP_BASE_PATH = http://localhost
// 	// 然后执行打包命令 npm run build。在打开下面4行注释
// 	// Router.LoadHTMLGlob("./dist/*.html") // npm打包成dist的路径
// 	// Router.Static("/favicon.ico", "./dist/favicon.ico")
// 	// Router.Static("/static", "./dist/assets")   // dist里面的静态资源
// 	// Router.StaticFile("/", "./dist/index.html") // 前端网页入口页面

// 	// Router.StaticFS(global.GIA_CONFIG.Local.Path, http.Dir(global.GIA_CONFIG.Local.StorePath)) // 为用户头像和文件提供静态地址
// 	// Router.Use(middleware.LoadTls())  // 如果需要使用https 请打开此中间件 然后前往 core/server.go 将启动模式 更变为 Router.RunTLS("端口","你的cre/pem文件","你的key文件")
// 	// 跨域，如需跨域可以打开下面的注释
// 	// Router.Use(middleware.Cors()) // 直接放行全部跨域请求
// 	//Router.Use(middleware.CorsByRules()) // 按照配置的规则放行跨域请求
// 	//global.GIA_LOG.Info("use middleware cors")

// 	// Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
// 	// global.GIA_LOG.Info("register swagger handler")
// 	// // 方便统一添加路由组前缀 多服务器上线使用

// 	PublicGroup := Router.Group("")
// 	{
// 		// 健康监测
// 		PublicGroup.GET("/health", func(c *gin.Context) {
// 			c.JSON(200, "ok")
// 		})
// 	}
// 	{
// 		fmt.Printf("systemRouter: %v\n", systemRouter)
// 		systemRouter.InitBaseRouter(PublicGroup) // 注册基础功能路由 不做鉴权
// 		// systemRouter.InitInitRouter(PublicGroup) // 自动初始化相关
// 	}
// 	PrivateGroup := Router.Group("")
// 	PrivateGroup.Use(middlewares.JWTAuth()).Use(middlewares.CasbinHandler())
// 	{
// 		// systemRouter.InitApiRouter(PrivateGroup)                 // 注册功能api路由
// 		// systemRouter.InitJwtRouter(PrivateGroup)                 // jwt相关路由
// 		// systemRouter.InitUserRouter(PrivateGroup)                // 注册用户路由
// 		// systemRouter.InitMenuRouter(PrivateGroup)                // 注册menu路由
// 		// systemRouter.InitSystemRouter(PrivateGroup)              // system相关路由
// 		// systemRouter.InitCasbinRouter(PrivateGroup)              // 权限相关路由
// 		// systemRouter.InitAutoCodeRouter(PrivateGroup)            // 创建自动化代码
// 		// systemRouter.InitAuthorityRouter(PrivateGroup)           // 注册角色路由
// 		// systemRouter.InitSysDictionaryRouter(PrivateGroup)       // 字典管理
// 		// systemRouter.InitAutoCodeHistoryRouter(PrivateGroup)     // 自动化代码历史
// 		// systemRouter.InitSysOperationRecordRouter(PrivateGroup)  // 操作记录
// 		// systemRouter.InitSysDictionaryDetailRouter(PrivateGroup) // 字典详情管理
// 		// systemRouter.InitAuthorityBtnRouterRouter(PrivateGroup)  // 字典详情管理

// 		// exampleRouter.InitExcelRouter(PrivateGroup)                 // 表格导入导出
// 		// exampleRouter.InitCustomerRouter(PrivateGroup)              // 客户路由
// 		// exampleRouter.InitFileUploadAndDownloadRouter(PrivateGroup) // 文件上传下载功能路由

// 		// Code generated by github.com/flipped-aurora/gin-vue-admin/server Begin; DO NOT EDIT.

// 		// Code generated by github.com/flipped-aurora/gin-vue-admin/server End; DO NOT EDIT.
// 	}

// 	// InstallPlugin(Router) // 安装插件

// 	// global.GIA_LOG.Info("router register success")
// 	log.Println("router register success")
// 	return Router
// }
