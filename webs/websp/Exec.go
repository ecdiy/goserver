package webs

import (
	"github.com/gin-gonic/gin"
	"github.com/cihub/seelog"
	"utils/gpa"
)

/**
Uri 规则
spPrefix sp名称前缀
 */

func Exec(g *gpa.Gpa, sp *Sp, ctx *gin.Context) (map[string]interface{}, int) {
	defer func() {
		if err := recover(); err != nil {
			delete(spCache, sp.Name)
		}
	}()
	params, code := sp.GetParams(ctx)
	if code == 200 {
		m, e := sp.Run(g, params...)
		if e != nil {
			seelog.Error("exec SP失败:", sp.Name)
			delete(spCache, sp.Name)
			return nil, 500
		}
		return m, 200
	} else {
		return nil, code
	}
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

func (sp *Sp) GetParams(ctx *gin.Context) ([]interface{}, int) {
	var params []interface{}
	for _, p := range sp.Params {
		vf, code := p.ValFunc(ctx, p)
		if code != 200 {
			seelog.Error("获取参数值出错：", p.ParamName)
			return nil, code
		}
		params = append(params, vf)
	}
	return params, 200
}
