package utils

import (
	"io"
	"log"
	"net"
	"net/http"
)

// 获取本机外网ip
func NetPublicIP() (ip string, err error) {
	// 发送GET请求到ifconfig.me
	resp, err := http.Get("http://ifconfig.me/ip")
	if err != nil {
		log.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	// body, err := ioutil.ReadAll(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error:", err)
		return
	}

	// 外网IP地址
	ip = string(body)
	// fmt.Println("External IP Address:", ip)

	return
}

// 获取本机内网ip
func NetLocalhostIP() string {
	// var ipV4, ipV6 string
	interfaceAddr, _ := net.InterfaceAddrs()
	for _, address := range interfaceAddr {
		ipNet, isVailIpNet := address.(*net.IPNet)
		if isVailIpNet && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				// ipV4 = ipNet.IP.String()
				return ipNet.IP.String()
			} else if ipNet.IP.To16() != nil {
				// ipV6 = ipNet.IP.String()
				return ipNet.IP.String()
			}
		}
	}
	// "127.0.0.1"
	return "localhost"
}
