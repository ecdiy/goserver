package webs

import (
	"utils/gpa"
	"github.com/cihub/seelog"
	"utils"
)

func SpCall(spName string, param *Param, result map[string]interface{},
	auth func(c *Param) *UserBase, enableCache bool) int {
	defer func() {
		if err := recover(); err != nil {
			delete(spCache, spName)
		}
	}()
	var sp *Sp
	var ext bool
	if utils.EnvIsDev {
		sp, ext = NewSpByName(Gpa, spName, auth)
	} else {
		sp, ext = spCache[spName]
		if !ext {
			sp, ext = NewSpByName(Gpa, spName, auth)
			if ext {
				spCache[spName] = sp
			}
		}
	}
	if !ext {
		return 404
	}
	params, code := sp.GetParams(param)
	if code == 200 {
		e := sp.Run(result, Gpa.Conn, params...)
		if e != nil {
			seelog.Error("exec SP失败:", sp.Name)
			delete(spCache, sp.Name)
			return 500
		}
		return 200
	} else {
		return code
	}
}

func SpExec(spName string, g *gpa.Gpa, param *Param, auth func(c *Param) *UserBase) int {
	return SpCall(spName, param, param.Out, auth, false)
}
