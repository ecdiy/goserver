package utils

import (
	"net"
	"strings"
)

func GetIp() string {
	return GetIpByExcludeName("VirtualBox")
}
func GetIpByExcludeName(ext ...string) string {
	ifs, e := net.Interfaces()
	if e != nil {
		return ""
	}
	for _, a := range ifs {
		ip, _ := a.Addrs()
		for _, ai := range ip {
			ax := ai.(*net.IPNet)
			if !ax.IP.IsLoopback() && strings.Index(ax.IP.String(), ".") > 0 {
				m := false
				for _, e := range ext {
					if strings.Index(a.Name, e) >= 0 {
						m = true
						break
					}
				}
				if !m {
					return ax.IP.String()
				}
			}
		}
	}
	return ""
}
