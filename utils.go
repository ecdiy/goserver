package utils

import (
	"net"
	"strings"
	"os"
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

func FmtImgDir(dir, uri string) (string, string) {
	ln := len(uri)
	if ln > 6 {
		x := uri[0:3] + "/" + uri[3:6]
		os.MkdirAll(dir+x, 0644)
		urx := x + "/" + uri[6:]
		path := dir + urx
		return path, urx
	}
	if ln > 3 {
		x := uri[0:3]
		os.MkdirAll(dir+x, 0644)
		urx := x + "/" + uri[3:]
		path := dir + urx
		return path, urx
	}
	return dir + uri, uri
}
