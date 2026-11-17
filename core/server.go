package core

import (
	"context"
	"fmt"
	"gmkdb/global"
	"gmkdb/initalice"
	"gmkdb/utils"
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type server interface {
	ListenAndServe() error
	ListenAndServeTLS(certFile, keyFile string) error
	Shutdown(ctx context.Context) error
}

func RunServer() {

	// 	// 初始化redis服务
	// if global.GIA_CONFIG.System.UseMultipoint || global.GIA_CONFIG.System.UseRedis {
	// 	initalice.Redis()
	// }

	// 从db加载jwt数据
	// if global.GIA_DB != nil {
	// 	system.LoadAll()
	// }

	// 加载路由
	Router := initalice.Routers()
	// 加载静态文件
	// Router.Static("/form-generator", "./resource/page")

	// 配置URL
	// address := fmt.Sprintf(":%d", global.GIA_CONFIG.System.Addr)
	address := fmt.Sprintf(":%d", global.GIA_CONFIG.System.Addr)
	if global.GIA_CONFIG.System.Addr == 80 || global.GIA_CONFIG.System.UseCertif {
		address = ""
	}

	// 启动服务
	// s := initServer(address, Router)
	s, rs := initServerPriority()(address, Router)

	wg := closeServe(s, rs)
	pid := strconv.FormatInt(int64(os.Getpid()), 10)
	url := getUrl(address)
	global.GIA_LOG.Info("Run server is ready:", zap.String("address", url))

	// 保证文本顺序输出
	// In order to ensure that the text order output can be deleted
	time.Sleep(10 * time.Microsecond)
	fmt.Printf(`
	server  run on : %s
	Swagger run on : %s/swagger/index.html
	Ping    run on : %s/health
	Program PID on : %s
`, url, url, url, pid)

	// global.GIA_LOG.Info("server run success on ", zap.String("address", address))
	// log.Println("server run success on ", address)

	// fmt.Println(`欢迎使用 gin-vue-admin当前版本:v2.5.3`, address)

	// global.GIA_LOG.Error(s.ListenAndServe().Error())
	// log.Println(s.ListenAndServe().Error())
	if global.GIA_CONFIG.System.UseCertif {
		go func() { // 同时启动 HTTP 重定向到 HTTPS
			global.GIA_LOG.Error(rs.ListenAndServe().Error())
		}()
		global.GIA_LOG.Error("TLS-" + s.ListenAndServeTLS(global.GIA_CONFIG.System.CertFile, global.GIA_CONFIG.System.KeyFile).Error())
	} else {
		global.GIA_LOG.Error(s.ListenAndServe().Error())
	}
	wg.Wait()
}

func getUrl(address string) string {
	// 获取ip
	ip, err := utils.NetPublicIP()
	if err != nil || !global.GIA_CONFIG.System.UseNetIP {
		ip = utils.NetLocalhostIP()
	}

	// 设置报头
	head := "http"
	if global.GIA_CONFIG.System.UseCertif {
		head = "https"
	}
	return fmt.Sprintf("%s://%s%s", head, ip, address)
}

func Run() {
	fmt.Println("Project Running...")
	// return time.Now()
}
