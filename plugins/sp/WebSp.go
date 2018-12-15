package sp

import (
	"github.com/ecdiy/goserver/gpa"
	"github.com/cihub/seelog"
	"strings"
	"regexp"
	"errors"
	"github.com/ecdiy/goserver/utils"
	"github.com/ecdiy/goserver/plugins"
	"runtime/debug"
)

type WebSp struct {
	Gpa          *gpa.Gpa
	SpParamDoMap map[string]ParamValFunc //存储过程参数处理规制
	SpCache      map[string]*Sp
}

func (ws *WebSp) Init(ele *utils.Element) {
	ws.Gpa = plugins.GetRef(ele, "Gpa").(*gpa.Gpa)
	ws.SpCache = make(map[string]*Sp)
	ws.SpParamDoMap = make(map[string]ParamValFunc)
	ws.SpParamDoMap["in"] = ParamIn
	ws.SpParamDoMap["ua"] = ParamUa
	pp, ppb := ele.AttrValue("ParamPrefix")
	if ppb {
		ppa := strings.Split(pp, ",")
		for _, name := range ppa {
			impl, ib := ParamDoMap[name]
			if ib {
				ws.SpParamDoMap[name] = impl
			} else {
				panic("Sp前缀没有实现")
			}
		}
	} else {
		for name := range ParamDoMap {
			ws.SpParamDoMap[name] = ParamDoMap[name]
		}
		seelog.Warn("没有设置ParamPrefix过滤规则，使用所有规则,", ele.ToString())
	}
}

func (ws *WebSp) NewSp(val []string) (*Sp, bool) {
	sp := &Sp{Name: val[0]}
	//if len(val) < 3 || len(strings.TrimSpace(val[2])) == 0 {
	//	seelog.Warn("没有返回值的参数申明")
	//	return sp, false
	//}
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
				spn, spe := ws.GetParam(pName, pType, sp)
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

func (ws *WebSp) NewSpByName(spName string) (*Sp, bool) {
	info, e := ws.Gpa.ListString(SqlSpInfo, spName)
	if e != nil || len(info) != 3 {
		seelog.Warn("存储过程不存在:", spName, e)
		return &Sp{}, false
	}
	return ws.NewSp(info)
}

func (ws *WebSp) GetRunSp(spName string) *Sp {
	var sp *Sp
	var ext bool
	if utils.EnvIsDev {
		sp, ext = ws.NewSpByName(spName)
	} else {
		sp, ext = ws.SpCache[spName]
		if !ext {
			sp, ext = ws.NewSpByName(spName)
			if ext {
				ws.SpCache[spName] = sp
			}
		}
	}
	if !ext {
		return nil
	}
	return sp
}

func (ws *WebSp) SpExec(spName string, param *utils.Param) int {
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
			seelog.Error("SP do fail: spName=", spName, ";param=", param)
			seelog.Errorf("\n\t%v", err)
			delete(ws.SpCache, spName)
		}
	}()
	sp := ws.GetRunSp(spName)
	if sp == nil {
		return 404
	}
	params, code := ws.GetParams(param, sp)
	if code == 200 {
		e := sp.Run(param.Out, ws.Gpa.Conn, params...)
		if e != nil {
			delete(ws.SpCache, sp.Name)
			return 500
		}
		return 200
	} else {
		return code
	}
}

func (ws *WebSp) GetParams(wb *utils.Param, sp *Sp) ([]interface{}, int) {
	var params []interface{}
	for _, p := range sp.Params {
		vf, code := p.ValFunc(wb, p)
		if code != 200 {
			return nil, code
		}
		params = append(params, vf)
	}
	return params, 200
}

func (ws *WebSp) GetParam(ParamName, pType string, sp *Sp) (*SpParam, error) {
	p := &SpParam{ParamName: ParamName}
	if pType == "bigint" || pType == "int" {
		p.DefaultVal = "0"
	} else {
		p.DefaultVal = ""
	}
	paraHead := ""
	for k, v := range ws.SpParamDoMap {
		paraHead += k + ","
		if strings.Index(ParamName, k) == 0 {
			p.ParamName = p.ParamName[len(k):]
			p.ValFunc = v
			return p, nil
		}
	}
	seelog.Error("合法参数以(", paraHead, "开头)未知参数格式，", ParamName)
	return p, errors.New("未知参数格式")
}
