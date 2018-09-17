package webs

import (
	"strings"
	"regexp"
	"github.com/cihub/seelog"
	"utils/gpa"
	"github.com/gin-gonic/gin"
)

func spInitCache(g *gpa.Gpa, auth func(c *gin.Context) (bool, int64), spPrefix string) {
	spReloadFun[spPrefix] = auth
	list, err := g.ListArrayString(SqlSpAll)
	if err != nil {
		seelog.Error("查询所有的存储过程出错：", err)
	} else {
		if list != nil {
			sps := ""
			for _, val := range list {
				if strings.LastIndex(val[0], spPrefix) == len(val[0])-len(spPrefix) {
					sp, b := NewSp(val, auth)
					if b {
						spCache [sp.Name] = sp
						sps += sp.Sql + ","
					}
				}
			}
			seelog.Info("~~~~", spPrefix, ";\n\t", sps)
		}
	}
}

func NewSp(val []string, auth func(c *gin.Context) (bool, int64)) (*Sp, bool) {
	sp := &Sp{Name: val[0], Auth: auth}
	if len(val) < 3 || len(strings.TrimSpace(val[2])) == 0 {
		seelog.Warn("没有返回值的参数申明")
		return sp, false
	}
	rowReg := regexp.MustCompile(`\r|\n`)
	rows := rowReg.Split(val[2], -1)
	if len(rows) < 1 {
		return sp, false
	}
	resInfo := strings.Split(rows[0], ",")
	for _, str := range resInfo {
		v := strings.Split(str, ":")
		spr := &SpResult{}
		if len(v) >= 2 {
			spr.Type = v[1]
		} else {
			spr.Type = "list"
		}
		if len(v) >= 1 {
			spr.Name = strings.TrimSpace(v[0])
			if len(spr.Name) > 0 {
				sp.Result = append(sp.Result, spr)
			}
		}
	}
	sp.Sql = "call " + sp.Name + "("
	if len(strings.TrimSpace(val[1])) > 1 {
		var spitBySpaceRegexp, _ = regexp.Compile(`\s+`)
		pNum := ""
		var spParams []*SpParam
		paramArray := strings.Split(val[1], ",")
		for _, p := range paramArray {
			if len(pNum) > 0 {
				pNum += ","
			}
			pNum += "?"
			pTrim := strings.TrimSpace(p)
			idxArray := spitBySpaceRegexp.FindIndex([]byte(pTrim))
			if idxArray != nil && idxArray[0] > 0 {
				pType := strings.ToLower(pTrim[idxArray[1]:])
				idx := strings.Index(pType, "(")
				if idx > 0 {
					pType = strings.TrimSpace(pType[0:idx])
				}
				pName := strings.TrimSpace(pTrim[0:idxArray[1]])
				spn, spe := sp.GetParam(pName, pType)
				if spe != nil {
					return sp, false
				}
				spParams = append(spParams, spn)
			}
		}
		sp.Params = spParams
		sp.Sql += pNum
	}
	sp.Sql += ")"
	return sp, true
}

func NewSpByName(g *gpa.Gpa, spName string, auth func(c *gin.Context) (bool, int64)) (*Sp, bool) {
	info, e := g.ListString(SqlSpInfo, spName)
	if e != nil || len(info) != 3 {
		seelog.Warn("存储过程不存在:", spName, e)
		return &Sp{}, false
	}
	return NewSp(info, auth)
}
