package initalice

import (
	"fmt"
	"gmkdb/global"
	"log"
	"net"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

func initdetectRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	detectRouter := Router.Group("")
	{
		// *
		// 获取系统基本设置
		detectRouter.GET("/config", func(c *gin.Context) {
			c.JSON(200, global.GIA_CONFIG)
		})

		// *
		// 测试status状态 (/status?status=200)
		detectRouter.GET("/status", func(c *gin.Context) {
			s := c.Query("status")
			status, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				log.Printf("status: %v\n", status)
				status = 200
			}

			c.JSON(int(status), gin.H{"status": status})
		})

		// *
		// 获取系统参数
		detectRouter.GET("/sysinfo", func(c *gin.Context) {
			pid := os.Getpid() // 直接调用函数
			v, _ := mem.VirtualMemory()
			percent, _ := cpu.Percent(time.Second, false)
			info, _ := host.Info()
			interfaces, _ := net.Interfaces()
			h := gin.H{}

			// 当前进程 pid
			h["pid"] = pid
			// fmt.Printf("当前进程的 PID 是: %d\n", pid)
			// 输出内存信息
			h["mem"] = gin.H{
				"Total":       v.Total / 1024 / 1024,
				"Used":        v.Used / 1024 / 1024,
				"Available":   v.Available / 1024 / 1024,
				"UsedPercent": v.UsedPercent,
			}
			// fmt.Printf("总内存: %v MB\n", v.Total/1024/1024)
			// fmt.Printf("已用内存: %v MB\n", v.Used/1024/1024)
			// fmt.Printf("可用内存: %v MB\n", v.Available/1024/1024)
			// fmt.Printf("内存使用率: %.2f%%\n", v.UsedPercent)

			// 获取硬盘分区信息
			part := gin.H{}
			partitions, _ := disk.Partitions(false)
			for _, partition := range partitions {
				usage, _ := disk.Usage(partition.Mountpoint)
				part[partition.Device] = gin.H{
					"Mountpoint":  partition.Mountpoint,
					"Fstype":      partition.Fstype,
					"Total":       usage.Total / 1024 / 1024 / 1024,
					"Used":        usage.Used / 1024 / 1024 / 1024,
					"Free":        usage.Free / 1024 / 1024 / 1024,
					"UsedPercent": usage.UsedPercent,
				}
				// fmt.Printf("设备: %s\n", partition.Device)
				// fmt.Printf("挂载点: %s\n", partition.Mountpoint)
				// fmt.Printf("文件系统类型: %s\n", partition.Fstype)
				// // 获取分区使用情况
				// fmt.Printf("总大小: %v GB\n", usage.Total/1024/1024/1024)
				// fmt.Printf("已用空间: %v GB\n", usage.Used/1024/1024/1024)
				// fmt.Printf("可用空间: %v GB\n", usage.Free/1024/1024/1024)
				// fmt.Printf("使用率: %.2f%%\n", usage.UsedPercent)
				// fmt.Println("------")
			}
			h["partitions"] = part

			// 获取 CPU 使用率
			h["cpu"] = percent[0]
			// fmt.Printf("CPU 使用率: %.2f%%\n", percent[0])

			// 获取主机信息
			h["info"] = gin.H{
				"Hostname":      info.Hostname,
				"OS":            info.OS,
				"Platform":      info.Platform,
				"KernelVersion": info.KernelVersion,
				"Uptime":        info.Uptime / 3600,
			}
			// fmt.Printf("主机名: %s\n", info.Hostname)
			// fmt.Printf("操作系统: %s\n", info.OS)
			// fmt.Printf("平台: %s\n", info.Platform)
			// fmt.Printf("内核版本: %s\n", info.KernelVersion)
			// fmt.Printf("运行时间: %v 小时\n", info.Uptime/3600)

			// 获取网络接口信息
			ifs := gin.H{}
			for _, iface := range interfaces {
				newIf := gin.H{
					"HardwareAddr": iface.HardwareAddr,
				}
				// fmt.Printf("接口名称: %s\n", iface.Name)
				// fmt.Printf("MAC 地址: %s\n", iface.HardwareAddr)

				if addrs, err := iface.Addrs(); err != nil {
					for i, addr := range addrs { // iface.Addrs
						// fmt.Printf("IP 地址: %s\n", addr.Addr)
						newIf["IP_"+strconv.FormatInt(int64(i), 10)] = addr.String()
					}

				}
				ifs[iface.Name] = newIf
			}
			h["iface"] = ifs

			c.JSON(200, h)
		})

		// *
		// 获取终止系统进程(自动获取本服务进程)
		detectRouter.GET("/exit", func(c *gin.Context) {
			pid := os.Getpid()     // 直接调用函数
			info, _ := host.Info() // 查看系统进行区分

			// if info.OS == "windows" { // "golang.org/x/sys/windows"
			// 	// // Windows: Generate a CTRL+C event
			// 	r1, _, err := windows.NewLazySystemDLL("kernel32.dll").NewProc("GenerateConsoleCtrlEvent").
			// 		Call(uintptr(windows.CTRL_C_EVENT), uintptr(pid))
			// 	if r1 == 0 {
			// 		log.Printf("failed to generate CTRL+C event for PID %d: %v\n", pid, err.Error())
			// 	}

			// 	// * 参考代码 // process.Kill(): 调用'Kill'方法杀死当前进程(等同于kill -9), 无法进行拦截
			// 	// // process, _:= os.FindProcess(pid)
			// } else
			if info.OS == "linux" { // "syscall"
				// // Linux: send SIGTERM sigal
				err := syscall.Kill(pid, syscall.SIGTERM) // 发送 SIGTERM 信号
				if err != nil {
					fmt.Println("终止进程失败:", err)
				}
			} else {
				process, err := os.FindProcess(pid)
				if err != nil {
					fmt.Printf("Failed to find process for PID %d: %v\n", pid, err)
				} else {
					// 使用 os.Process.Signal 发送 syscall.SIGTERM 信号 //
					err := process.Signal(syscall.SIGTERM)
					if err != nil {
						log.Printf("Failed to send SIGTERM: %v\n", err)
					}
				}
			}

			c.JSON(200, fmt.Sprintf("PID: %d", pid))
		})
	}
	return detectRouter
}
