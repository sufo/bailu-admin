/**
 * Create by sufo
 * @Email ouamour@Gmail.com
 *
 * @Desc 客户端工具
 */

package utils

import (
	"errors"
	"fmt"
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
	"github.com/mssola/user_agent"
	"bailu/global/consts"
	"bailu/pkg/ip2region"
	"net"
	"net/http"
	"os"
	"strings"
)

func GetIP(r *http.Request) (string, error) {
	ip := r.Header.Get("X-Real-IP")
	if net.ParseIP(ip) != nil {
		return ip, nil
	}

	ip = r.Header.Get("X-Forward-For")
	for _, i := range strings.Split(ip, ",") {
		if net.ParseIP(i) != nil {
			return i, nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	if net.ParseIP(ip) != nil {
		return ip, nil
	}
	return "", errors.New("no valid ip found")
}

func GetBrowser(r *http.Request) (string, string) {
	ua := user_agent.New(r.UserAgent())
	return ua.Engine()
	//fmt.Printf("%v\n", name) // => "AppleWebKit"
	//fmt.Printf("%v\n", version)
}

func GetOs(ua *user_agent.UserAgent) string {
	return ua.OS()
}

// LocalIP get the host machine local IP address
// 不准确
func GetLocalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if isPrivateIP(ip) {
				return ip.String(), nil
			}
		}
	}
	return "", fmt.Errorf("IP not found!")
}

func isPrivateIP(ip net.IP) bool {
	var privateIPBlocks []*net.IPNet
	for _, cidr := range []string{
		// don't check loopback ips
		//"127.0.0.0/8",    // IPv4 loopback
		//"::1/128",        // IPv6 loopback
		//"fe80::/10",      // IPv6 link-local
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
	} {
		_, block, _ := net.ParseCIDR(cidr)
		privateIPBlocks = append(privateIPBlocks, block)
	}

	for _, block := range privateIPBlocks {
		if block.Contains(ip) {
			return true
		}
	}
	return false
}

// 通过ip获取address
func GetAddr(ip string) string {
	if ip == "" || ip2region.VIndex == nil {
		return ""
	}
	// 2、用全局的 vIndex 创建带 VectorIndex 缓存的查询对象。
	searcher, err := xdb.NewWithVectorIndex(ip2region.DbPath, ip2region.VIndex)
	if err != nil {
		fmt.Printf("failed to create searcher with vector index: %s\n", err)
		return ""
	}
	defer searcher.Close()
	region, err := searcher.SearchByStr(ip)
	if err != nil {
		fmt.Printf("failed to SearchIP(%s): %s\n", ip, err)
		return ""
	}
	return region
}

func GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		fmt.Println("get hostname failed, err = ", err.Error())
		return ""
	}
	return hostname
}

func IsMobile(userAgent string) bool {
	ua := user_agent.New(userAgent)
	return ua.Mobile()
}

// PC or mobile
func DeviceType(userAgent string) string {
	ua := user_agent.New(userAgent)
	if ua.Mobile() {
		return consts.DEVICE_MOBILE
	} else {
		return consts.DEVICE_PC
	}
}
