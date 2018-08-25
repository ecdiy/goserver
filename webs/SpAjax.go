package webs

import (
	"github.com/gin-gonic/gin"
	"github.com/cihub/seelog"

	"utils/gpa"
	"fmt"
	"encoding/json"
	"strings"
	"regexp"
	"strconv"
	"errors"
	"utils"
)

var spCache = make(map[string]Sp)

var simpleNameReg = regexp.MustCompile(`(^\W*)|([.].*$)`)
/**
URL 映射到存储过程调用，返回json数据格式
SpAjax(true,"/sp/","Sp",authFun)
 */
func SpAjax(uri string, g *gpa.Gpa, eng *gin.Engine, spPrefix string, auth func(c *gin.Context) bool) {
	eng.GET(uri+"Reload", func(c *gin.Context) {
		seelog.Info("Reload Sp Cache")
		spCache = make(map[string]Sp)
		spInitCache(g)
		utils.OK.OutJSON(c, nil)
	})
	if !gin.IsDebugging() {
		spInitCache(g)
	}
	if auth != nil {
		eng.POST(uri+"/*sp", func(c *gin.Context) {
			if auth(c) {
				spName := spPrefix + simpleNameReg.ReplaceAllString(c.Param("sp"), "")
				var sp Sp
				if gin.IsDebugging() {
					sp, _ = LoadSpFromDb(g, spName)
				} else {
					sp = spCache[spName]
				}
				data, err := GinSp(g, sp, c)
				if err == nil {
					utils.OK.OutJSON(c, data)
				} else {
					seelog.Error("数据存储过程错误:"+spName, "\n\t", err)
					utils.StErrorDb.OutJSON(c, nil)
				}
			} else {
				utils.StErrorToken.OutJSON(c, nil)
			}
		})
	}
}

func spInitCache(g *gpa.Gpa) {
	list, err := g.ListArrayString(SqlSpAll)
	if err != nil {
		seelog.Error("查询所有的存储过程出错：", err)
	} else {
		if list != nil {
			sps := ""
			for _, val := range list {
				sp, b := newSp(val)
				if b {
					spCache [sp.Name] = sp
					sps += sp.Sql + ","
				}
			}
			seelog.Info("~~sp:~~", sps)
		}
	}
}

func newSp(val []string) (Sp, bool) {
	sp := Sp{Name: val[0]}
	if len(val) < 3 {
		return sp, false
	}

	if len(val[2]) > 0 {
		sp.Info = make(map[string]interface{})
		err := json.Unmarshal([]byte(val[2]), &sp.Info)
		if err != nil {
			seelog.Error(sp.Name, " json format error:", err)
		}
	}
	sp.Sql = "call " + sp.Name + "("
	if len(strings.TrimSpace(val[1])) > 1 {
		var spitBySpaceRegexp, _ = regexp.Compile(`\s+`)
		pNum := ""
		var spParams []SpParam
		paramArray := strings.Split(val[1], ",")
		for _, p := range paramArray {
			if len(pNum) > 0 {
				pNum += ","
			}
			pNum += "?"
			pTrim := strings.TrimSpace(p)
			idxArray := spitBySpaceRegexp.FindIndex([]byte(pTrim))
			if idxArray != nil && idxArray[0] > 0 {
				length := 0
				pType := strings.ToLower(pTrim[idxArray[1]:])
				idx := strings.Index(pType, "(")
				if idx > 0 {
					length, _ = strconv.Atoi(pType[idx+1 : len(pType)-1])
					pType = strings.TrimSpace(pType[0:idx])
				}
				pName := strings.TrimSpace(pTrim[0:idxArray[1]])
				spParam := SpParam{ParamName: pName, ParamType: pType, Length: length}
				eSp := spParam.Init()
				if eSp != nil {
					return sp, false
				}
				spParams = append(spParams, spParam)
			}
		}
		sp.Params = spParams
		sp.Sql += pNum
	}
	sp.Sql += ")"
	return sp, true
}

func (c *SpParam) GinParam(ctx *gin.Context, param map[string]interface{}) (string, error) {
	v, b := ctx.Get(c.ParamName)
	if b {
		return fmt.Sprint(v), nil
	} else {
		seelog.Error("ctx.Get not find.", c.ParamName)
		return "", errors.New("not find[" + c.ParamName + "]")
	}
}

func (c *SpParam) Init() error {
	//	if strings.Index(c.ParamName, "session2") == 0 {
	//		c.ParamName = c.ParamName[8:]
	//		if c.ParamName[0:1] == "_" {
	//			c.ParamName = c.ParamName[1:]
	//		}
	//		c.Param = strings.Split(c.ParamName, "_")
	//		if len(c.Param) > 2 {
	//			seelog.Error("错误的存储过程参数:", c.ParamName, ",不支持多个[_],只能是一个")
	//		}
	//		sp.SessionName = c.ParamName
	//		c.ValFunc = c.Session2
	//		c.FuncType = "session"
	//	} else
	if strings.Index(c.ParamName, "gin") == 0 {
		c.ParamName = c.ParamName[3:]
		c.ValFunc = c.GinParam
		//	c.FuncType = "session"
		return nil
	}
	if strings.Index(c.ParamName, "in") == 0 {
		//c.FuncType = "context"
		c.ParamName = c.ParamName[2:]
		c.ValFunc = c.InParam
		return nil
	}
	//}
	//func (c *SpParam) Session2(ctx *gin.Context, defaultVal map[string]map[string]string) (string, error) {
	//	v, st := c.Session(ctx, defaultVal)
	//	if st.Code == utils.ST_OK.Code {
	//		return v, st
	//	}
	//	return "0", nil
	seelog.Error("合法参数以(in,gin开头)未知参数格式，", c.ParamName)
	return errors.New("未知参数格式")
}

func (c *SpParam) InParam(ctx *gin.Context, param map[string]interface{}) (string, error) {
	if param != nil {
		v, b := param[c.ParamName]
		if b {
			return fmt.Sprint(v), nil
		}
	}
	v := ctx.Param(c.ParamName)
	if v == "" {
		v = ctx.PostForm(c.ParamName)
	}
	if v == "" {
		v = ctx.Query(c.ParamName)
	}
	return v, nil
}

func LoadSpFromDb(g *gpa.Gpa, spName string) (Sp, error) {
	info, e := g.ListString("select name,CONVERT(param_list USING utf8) param_list,`comment` from mysql.proc c where db=DATABASE() and `type`='PROCEDURE' and name=?", spName)
	if e != nil || len(info) == 0 {
		return Sp{}, e
	}
	if len(info) < 3 {
		seelog.Error("存储过程不存在:", spName)
		return Sp{}, errors.New("sp not exist")
	}
	sp, sb := newSp(info)
	if sb {
		return sp, nil
	} else {
		return sp, errors.New("parameter name error")
	}
}

/**
Uri 规则
spPrefix sp名称前缀
 */

func GinSp(g *gpa.Gpa, sp Sp, ctx *gin.Context) (map[string][]map[string]string, error) {
	defer func() {
		if err := recover(); err != nil {
			delete(spCache, sp.Name)
		}
	}()
	var params []interface{}
	if len(sp.Params) > 0 {
		data, je := Parameter(ctx)
		if je != nil {
			seelog.Error("RawData JSON error:", je)
		}

		for _, p := range sp.Params {
			vf, errVf := p.ValFunc(ctx, data.Param)
			if errVf != nil {
				seelog.Error("获取参数值出错：", p.ParamName)
				return nil, errVf
			}
			params = append(params, vf)
		}
	}
	m, e := g.Sp(sp.Sql, params...)
	if e != nil {
		seelog.Error("exec SP失败:", sp.Name)
		delete(spCache, sp.Name)
	}
	return m, e
}
