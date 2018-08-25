package webs

import (
	"github.com/gin-gonic/gin"
	"github.com/cihub/seelog"
	"utils/gpa"
	"strings"
	"errors"
)

//--sp
type Sp struct {
	Sql, Name, SessionName string
	Params                 []*SpParam
	Result                 []*SpResult
	Auth                   func(c *gin.Context) (bool, int64)
}

type SpResult struct {
	Name, Type string //[ [total,object],[list,list] ]
}
type SpParam struct {
	ParamName string
	ValFunc   ParamValFunc
}

type ParamValFunc func(ctx *WebBase, p *SpParam) (interface{}, int)

var spCache = make(map[string]*Sp)

//--
func (sp *Sp) InParam(ctx *WebBase, p *SpParam) (interface{}, int) {
	v := ctx.String(p.ParamName)
	return v, 200
}

func (sp *Sp) GinParam(ctx *WebBase, p *SpParam) (interface{}, int) {
	v, b := ctx.Context.Get(p.ParamName)
	if b {
		return v, 200
	} else {
		if sp.Auth != nil {
			auth, _ := sp.Auth(ctx.Context)
			if !auth {
				return "", 401
			}
			v, b := ctx.Context.Get(p.ParamName)
			if b {
				return v, 200
			} else {
				return "", 404
			}
		}
		seelog.Error("ctx.Get not find.", p.ParamName)
		return "", 401
	}
}

//--
func (sp *Sp) GetParams(wb  *WebBase) ([]interface{}, int) {
	var params []interface{} 
	for _, p := range sp.Params {
		vf, code := p.ValFunc(wb, p)
		if code != 200 {
			seelog.Error("获取参数值出错：", p.ParamName)
			return nil, code
		}
		params = append(params, vf)
	}
	return params, 200
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

func (sp *Sp) Run(g *gpa.Gpa, params ...interface{}) (map[string]interface{}, error) {
	data := make(map[string]interface{})
	rows, err := g.Conn.Query(sp.Sql, params...)
	defer rows.Close()
	if err != nil {
		seelog.Error("调用存储过程出错了.", sp.Sql, params, "\n\t", err)
		return data, err
	}
	for node := 0; node < len(sp.Result); node++ {
		cols, err := rows.Columns()
		if err != nil {
			seelog.Error("获取结果集columns出错:", err)
			return data, err
		}
		r := sp.Result[node]
		if r.Type == "list" {
			var list []map[string]string
			for rows.Next() {
				list = append(list, gpa.RowToMap(rows, cols))
			}
			data[sp.Result[node].Name] = list
		}
		if r.Type == "object" {
			data[sp.Result[node].Name] = gpa.RowToMap(rows, cols)
		}
		if !rows.NextResultSet() {
			break
		}
	}
	return data, nil
}
