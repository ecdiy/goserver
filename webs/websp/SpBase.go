package webs

import (
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"strings"
	"utils/gpa"
	"errors"
)

var spCache = make(map[string]*Sp)

//var simpleNameReg = regexp.MustCompile(`(^\W*)|([.].*$)`)

func spInitCache(g *gpa.Gpa) {
	list, err := g.ListArrayString(SqlSpAll)
	if err != nil {
		seelog.Error("查询所有的存储过程出错：", err)
	} else {
		if list != nil {
			sps := ""
			for _, val := range list {
				sp, b := NewSp(val)
				if b {
					spCache [sp.Name] = sp
					sps += sp.Sql + ","
				}
			}
			seelog.Info("~~sp:~~", sps)
		}
	}
}

func (sp *Sp) GetParam(ParamName string) (*SpParam, error) {

	p := &SpParam{ParamName: ParamName}
	if strings.Index(ParamName, "gin") == 0 {
		p.ParamName = p.ParamName[3:]
		p.ValFunc = sp.GinParam
		return p, nil
	}
	if strings.Index(ParamName, "in") == 0 {
		p.ParamName = p.ParamName[2:]
		p.ValFunc = sp.InParam
		return p, nil
	}

	seelog.Error("合法参数以(in,gin开头)未知参数格式，", ParamName)
	return p, errors.New("未知参数格式")
}

func (sp *Sp) InParam(ctx *gin.Context, p *SpParam) (interface{}, int) {
	v := ctx.Param(p.ParamName)
	if v == "" {
		v = ctx.PostForm(p.ParamName)
	}
	if v == "" {
		v = ctx.Query(p.ParamName)
	}
	return v, 200
}

func (sp *Sp) GinParam(ctx *gin.Context, p *SpParam) (interface{}, int) {
	v, b := ctx.Get(p.ParamName)
	if b {
		return v, 200
	} else {
		seelog.Error("ctx.Get not find.", p.ParamName)
		return "", 404
	}
}
