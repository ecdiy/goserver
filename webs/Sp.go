package webs

import (
	"github.com/cihub/seelog"
	"utils/gpa"
	"strings"
	"errors"
	"database/sql"
)

//--sp
type Sp struct {
	Sql, Name, SessionName string
	Params                 []*SpParam
	Result                 []*SpResult
	Auth                   func(c *Param) (bool, int64)
}

type SpResult struct {
	Name, Type string //[ [total,object],[list,list] ]
}
type SpParam struct {
	ParamName, DefaultVal string
	ValFunc               ParamValFunc
}

type ParamValFunc func(ctx *Param, p *SpParam) (interface{}, int)

//--
func (sp *Sp) UaParam(wb *Param, p *SpParam) (interface{}, int) {
	return wb.Ua, 200
}
func (sp *Sp) InParam(ctx *Param, p *SpParam) (interface{}, int) {
	v := ctx.String(p.ParamName)
	if v == "" {
		v = p.DefaultVal
	}
	return v, 200
}

func (sp *Sp) GkParam(ctx *Param, p *SpParam) (interface{}, int) {
	_, CallAuth := ctx.Context.Get("CallAuth")
	if !CallAuth && sp.Auth != nil {
		sp.Auth(ctx)
	}
	v, b := ctx.Context.Get(p.ParamName)
	if b {
		return v, 200
	} else {
		seelog.Error("ctx.Get not find.", p.ParamName)
		return p.DefaultVal, 200
	}
}

func (sp *Sp) GinParam(ctx *Param, p *SpParam) (interface{}, int) {
	v, b := ctx.Context.Get(p.ParamName)
	if b {
		return v, 200
	} else {
		_, CallAuth := ctx.Context.Get("CallAuth")
		if !CallAuth && sp.Auth != nil {
			auth, _ := sp.Auth(ctx)
			if !auth {
				seelog.Error("Gin获取参数值出错：SpName=", sp.Name, ";ParamName=", p.ParamName)
				return "", 401
			}
			v, b := ctx.Context.Get(p.ParamName)
			if b {
				return v, 200
			} else {
				seelog.Error("获取Gin参数错误:", p.ParamName)
				return "", 401
			}
		}
		seelog.Error("ctx.Get not find.", p.ParamName)
		return "", 401
	}
}

//--
func (sp *Sp) GetParams(wb *Param) ([]interface{}, int) {
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

func (sp *Sp) GetParam(ParamName, pType string) (*SpParam, error) {
	p := &SpParam{ParamName: ParamName}
	if pType == "bigint" || pType == "int" {
		p.DefaultVal = "0"
	} else {
		p.DefaultVal = ""
	}
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
	if strings.Index(ParamName, "ua") == 0 {
		p.ValFunc = sp.UaParam
		return p, nil
	}
	if strings.Index(ParamName, "gk") == 0 {
		p.ParamName = p.ParamName[2:]
		p.ValFunc = sp.GkParam
		seelog.Info("Type=", pType, ";default=", p.DefaultVal)
		return p, nil
	}
	seelog.Error("合法参数以(in,gin,ua,gk开头)未知参数格式，", ParamName)
	return p, errors.New("未知参数格式")
}

func (sp *Sp) Run(data map[string]interface{}, Conn *sql.DB, params ...interface{}) error {
	rows, err := Conn.Query(sp.Sql, params...)
	defer rows.Close()
	if err != nil {
		seelog.Error("调用存储过程出错了.", sp.Sql, params, "\n\t", err)
		return err
	}
	for node := 0; node < len(sp.Result); node++ {
		cols, err := rows.Columns()
		if err != nil {
			seelog.Error("获取结果集columns出错:", err)
			return err
		}
		r := sp.Result[node]
		if r.Type == "list" ||  r.Type == "l"  {
			var list []map[string]string
			for rows.Next() {
				list = append(list, gpa.RowToMap(rows, cols))
			}
			data[sp.Result[node].Name] = list
		} else if r.Type == "object" ||  r.Type == "o" {
			if rows.Next() {
				data[sp.Result[node].Name] = gpa.RowToMap(rows, cols)
			}
		} else {
			seelog.Warn("未知类型:", r.Type)
		}
		if !rows.NextResultSet() {
			break
		}
	}
	return nil
}
