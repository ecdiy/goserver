package webs

import (
	"github.com/cihub/seelog"
	"goserver/gpa"
	"database/sql"
)

//--sp
type Sp struct {
	Sql, Name, SessionName string
	Params                 []*SpParam
	Result                 []*SpResult
	//Auth                   func(c *Param) *UserBase
}

type SpResult struct {
	Name, Type string //[ [total,object],[list,list] ]
}
type SpParam struct {
	ParamName, DefaultVal string
	ValFunc               ParamValFunc
}

type ParamValFunc func(ctx *Param, p *SpParam) (interface{}, int)

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
		if r.Type == "list" || r.Type == "l" {
			var list []map[string]string
			for rows.Next() {
				list = append(list, gpa.RowToMap(rows, cols))
			}
			data[sp.Result[node].Name] = list
		} else if r.Type == "object" || r.Type == "o" {
			if rows.Next() {
				data[sp.Result[node].Name] = gpa.RowToMap(rows, cols)
			} else {
				data[sp.Result[node].Name] = make(map[string]string)
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
