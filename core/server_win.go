package core

import (
	"context"
	"crypto/tls"
	"fmt"
	"gmkdb/global"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// func initServer(address string, router *gin.Engine) server {
// 	return &http.Server{
// 		Addr:           address,
// 		Handler:        router,
// 		ReadTimeout:    20 * time.Second,
// 		WriteTimeout:   20 * time.Second,
// 		MaxHeaderBytes: 1 << 20,
// 	}
// }

type serverFunc func(address string, router *gin.Engine) (server, server)

func initServer(address string, router *gin.Engine) (server, server) {
	return &http.Server{
		Addr:           address,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 0, // 1 << 20,
	}, nil
}

func initServerTLS(address string, router *gin.Engine) (server, server) {
	return &http.Server{
		Addr:    address,
		Handler: router,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12, // 最低 TLS 版本(推荐 TLS 1.2+)
			// CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256}, // 安全的椭圆曲线
			// CipherSuites:     []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384, tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384},
		},
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		IdleTimeout:    0,
		MaxHeaderBytes: 1 << 20, // 1MB
	}, &http.Server{
		Addr:    ":80",
		Handler: redirectToHTTPS(router), // 自定义跳转逻辑
	}
}

// var certManager autocert.Manager

func initServerPriority() serverFunc {
	if global.GIA_CONFIG.System.UseCertif {
		return initServerTLS
	}
	return initServer
}

// 将 HTTP 请求重定向到 HTTPS (如 http://example.com->https://example.com)
func redirectToHTTPS(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.TLS == nil { // 非 HTTPS 请求
			fmt.Println(r.URL.Path, " -*- ", r.RequestURI)
			httpsURL := "https://" + r.Host + r.URL.Path
			http.Redirect(w, r, httpsURL, http.StatusMovedPermanently)
			return
		}
		handler.ServeHTTP(w, r) // 直接处理 HTTPS 请求
	})
}

func closeServe(s, rs server) *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)

	start := time.Now()
	// 捕获 SIGINT 信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	global.GIA_LOG.Debug("Core Activate end and Signal capture.")
	go func() {
		defer wg.Done()
		x := <-sigChan
		// // interrupt terminated
		global.GIA_LOG.Warn("End Close Server, exiting...", zap.String("signal", x.String()))

		// 关闭, 清理系统资源
		global.Shutdown(start)

		// 优雅关闭 HTTP 服务器(等待现有请求完成)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		if rs != nil {
			rs.Shutdown(ctx)
		}
		s.Shutdown(ctx)
		cancel()
		fmt.Println(" * System Program Closed.")
	}()

	return &wg
}
