package utils

import (
	"github.com/cihub/seelog"
	"os"
	"strconv"
	"strings"
	"github.com/gin-gonic/gin"
	"regexp"
)

//--所有常量
const (
	KeyBindAddr = "BindAddr"
	EnvProd     = "prod"
	EnvDev      = "dev"
)

var (
	EnvIsDev = true
)

func init() {
	logger, _ := seelog.LoggerFromConfigAsBytes([]byte(`
				<seelog type="sync" minlevel="info" maxlevel="error">
					<outputs formatid="main"><console/></outputs>
					<formats><format id="main" format="[%Level] %Date %Time [%File %Func %Line] %Msg%n"/></formats>
				</seelog>`))
	seelog.ReplaceLogger(logger)
	for _, kv := range os.Args {
		v := strings.TrimSpace(kv)
		if len(v) > 0 {
			vv := strings.Split(v, ";")
			for _, v2 := range vv {
				idx := strings.Index(v2, "=")
				if idx > 0 {
					name := v[0:idx]
					if name == "profile" {
						profile = v[idx+1:]
					} else {
						params[name] = v[idx+1:]
					}
				}
			}
		}
	}
	baseDev()
	if EnvIsDev {
		ip := GetIp()
		EnvParamSet("ImgHost", "http://"+ip)
	}
}

func baseDev() {
	if profile == "" {
		profile = "dev"
	}
	EnvIsDev = profile == EnvDev

	if EnvIsDev {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
}

var profile string
var params = make(map[string]string)

func EnvParam(key string) string {
	pv, pb := params[key]
	if pb {
		return pv
	}
	av := os.Getenv(key)
	if av != "" {
		return av
	}
	return ""
}

func FmtVal(v string) string {
	vv := v
	re := regexp.MustCompile("[$][{]([^}]*)[}]")
	gs := re.FindAllStringSubmatch(v, -1)
	if gs != nil {
		for _, p := range gs {
			vv = strings.Replace(vv, "${"+p[1]+"}", EnvParam(p[1]), -1)
		}
	}
	return vv
}

func EnvParamInt(key string, defaultVal int) int {
	v := EnvParam(key)
	if v != "" {
		iv, e := strconv.Atoi(v)
		if e == nil {
			return iv
		}
		return defaultVal
	}
	return defaultVal
}

func EnvParamSet(key, val string) {
	if key == "profile" {
		profile = val
		baseDev()
	} else {
		_, b := params[key]
		if !b {
			params[key] = val
		}
	}
}

/**
示例("dev","prod",`
Key1=111
Key2=abc
`)
*/
func ParamInit(prof, conf string) {
	if prof == profile {
		args := strings.Split(conf, "\n")
		for _, arg := range args {
			idx2 := strings.Index(arg, "#")
			if idx2 == 0 {
				continue
			}
			idx := strings.Index(arg, "=")
			if idx > 0 {
				n := arg[0:idx]
				_, b := params[n]
				if b {
					continue
				} else {
					av := os.Getenv(n)
					if av != "" {
						params[n] = strings.TrimSpace(av)
					} else {
						params[n] = strings.TrimSpace(arg[idx+1:])
					}
				}
			}
		}
		seelog.Info("profile=", prof, params)
	}
}
